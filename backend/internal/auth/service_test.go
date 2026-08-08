package auth

import (
	"context"
	"errors"
	"testing"
)

type mockRepo struct {
	users map[string]*userRecord
}

func (m *mockRepo) FindByEmail(ctx context.Context, email string) (*userRecord, error) {
	if u, exists := m.users[email]; exists {
		return u, nil
	}
	return nil, ErrNotFound
}

func (m *mockRepo) FindByID(ctx context.Context, id string) (*userRecord, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockRepo) Create(ctx context.Context, req RegisterRequest, passwordHash string) (string, error) {
	return "test-id", nil
}

func (m *mockRepo) CreatePendingAnggota(ctx context.Context, req SignupAnggotaRequest, passwordHash string) (string, error) {
	gID := req.GoogleID
	rec := &userRecord{
		ID:          "pending-id",
		Email:       req.Email,
		NamaLengkap: req.NamaLengkap,
		NoHp:        req.NoHp,
		RoleID:      4,
		RoleName:    "Anggota",
		Scope:       "anggota",
		Status:      "pending",
		GoogleID:    &gID,
	}
	m.users[req.Email] = rec
	return "pending-id", nil
}

func (m *mockRepo) CreateGoogleUser(ctx context.Context, req GoogleLoginRequest) (*userRecord, error) {
	gID := req.GoogleID
	if gID == "" {
		gID = "goog_" + req.Email
	}
	rec := &userRecord{
		ID:          "g-id-1",
		Email:       req.Email,
		NamaLengkap: req.NamaLengkap,
		NoHp:        req.NoHp,
		RoleID:      4,
		RoleName:    "Anggota",
		Scope:       "anggota",
		Status:      "pending",
		GoogleID:    &gID,
	}
	m.users[req.Email] = rec
	return rec, nil
}

func (m *mockRepo) UpdatePassword(ctx context.Context, userID, newHash string) error {
	return nil
}

func (m *mockRepo) UpdateGoogleID(ctx context.Context, userID, googleID string) error {
	for _, u := range m.users {
		if u.ID == userID {
			u.GoogleID = &googleID
			return nil
		}
	}
	return ErrNotFound
}

func TestGoogleLogin_NewUserAutoCreate(t *testing.T) {
	repo := &mockRepo{users: make(map[string]*userRecord)}
	svc := NewService(repo, "secret123")

	// 1. Without NoHp/UnitID -> expect ErrUserNotFound (prompts profile completion)
	resp, err := svc.GoogleLogin(context.Background(), GoogleLoginRequest{
		Email:       "test.google@gmail.com",
		NamaLengkap: "Test Google User",
		GoogleID:    "goog_test123",
	})

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}

	if resp != nil {
		t.Fatalf("expected nil response for unregistered user")
	}

	// 2. With NoHp/UnitID -> creates pending user and expects ErrUserInactive
	resp, err = svc.GoogleLogin(context.Background(), GoogleLoginRequest{
		Email:       "test.google@gmail.com",
		NamaLengkap: "Test Google User",
		GoogleID:    "goog_test123",
		NoHp:        "081234567890",
		UnitID:      "unit-1",
	})

	if !errors.Is(err, ErrUserInactive) {
		t.Fatalf("expected ErrUserInactive, got: %v", err)
	}
}


func TestSignupAnggota_WithGoogleID(t *testing.T) {
	repo := &mockRepo{users: make(map[string]*userRecord)}
	svc := NewService(repo, "secret123")

	id, err := svc.SignupAnggota(context.Background(), SignupAnggotaRequest{
		Email:       "new.member@gmail.com",
		Password:    "password123",
		NamaLengkap: "New Member",
		NoHp:        "081299998888",
		UnitID:      "unit-1",
		Tingkatan:   "Pra Dasar",
		GoogleID:    "goog_newmember",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if id != "pending-id" {
		t.Errorf("expected id pending-id, got %s", id)
	}

	userRec := repo.users["new.member@gmail.com"]
	if userRec.GoogleID == nil || *userRec.GoogleID != "goog_newmember" {
		t.Errorf("expected GoogleID goog_newmember, got %v", userRec.GoogleID)
	}
}
