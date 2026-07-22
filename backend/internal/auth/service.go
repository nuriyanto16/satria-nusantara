package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"satria-nusantara/backend/pkg/middleware"
)

// Sentinel errors for domain-level error handling.
var (
	ErrNotFound         = errors.New("user not found")
	ErrInvalidPassword  = errors.New("email atau password salah")
	ErrEmailExists      = errors.New("email sudah terdaftar")
	ErrUserInactive     = errors.New("akun tidak aktif, hubungi administrator")
)

// Service defines the business logic contract for authentication.
type Service interface {
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	GoogleLogin(ctx context.Context, req GoogleLoginRequest) (*AuthResponse, error)
	GetMe(ctx context.Context, userID string) (*User, error)
	Register(ctx context.Context, req RegisterRequest) (string, error)
	ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error
}

// service is the concrete implementation of Service.
type service struct {
	repo      Repository
	jwtSecret string
}

// NewService creates a new auth service.
func NewService(repo Repository, jwtSecret string) Service {
	return &service{repo: repo, jwtSecret: jwtSecret}
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	rec, err := s.repo.FindByEmail(ctx, req.Email)
	if errors.Is(err, ErrNotFound) {
		return nil, ErrInvalidPassword // Don't leak user existence
	}
	if err != nil {
		return nil, err
	}

	if rec.Status != "aktif" {
		return nil, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(rec.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	cabangID := ""
	if rec.CabangID != nil {
		cabangID = *rec.CabangID
	}
	unitID := ""
	if rec.UnitID != nil {
		unitID = *rec.UnitID
	}

	token, err := middleware.GenerateToken(
		rec.ID, rec.Email, rec.Scope, cabangID, unitID, s.jwtSecret, rec.RoleID,
	)
	if err != nil {
		return nil, err
	}

	user := recordToUser(rec)
	return &AuthResponse{Token: token, User: user}, nil
}

func (s *service) GoogleLogin(ctx context.Context, req GoogleLoginRequest) (*AuthResponse, error) {
	if req.Email == "" {
		return nil, errors.New("email wajib diisi")
	}

	rec, err := s.repo.FindByEmail(ctx, req.Email)
	if errors.Is(err, ErrNotFound) {
		newRec, createErr := s.repo.CreateGoogleUser(ctx, req)
		if createErr != nil {
			return nil, createErr
		}
		rec = newRec
	} else if err != nil {
		return nil, err
	}

	if rec.Status != "aktif" {
		return nil, ErrUserInactive
	}

	if req.GoogleID != "" && (rec.GoogleID == nil || *rec.GoogleID == "") {
		_ = s.repo.UpdateGoogleID(ctx, rec.ID, req.GoogleID)
	}

	cabangID := ""
	if rec.CabangID != nil {
		cabangID = *rec.CabangID
	}
	unitID := ""
	if rec.UnitID != nil {
		unitID = *rec.UnitID
	}

	token, err := middleware.GenerateToken(
		rec.ID, rec.Email, rec.Scope, cabangID, unitID, s.jwtSecret, rec.RoleID,
	)
	if err != nil {
		return nil, err
	}

	user := recordToUser(rec)
	return &AuthResponse{Token: token, User: user}, nil
}


func (s *service) GetMe(ctx context.Context, userID string) (*User, error) {
	rec, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return recordToUser(rec), nil
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (string, error) {
	// Check if email already exists
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return "", err
	}
	if existing != nil {
		return "", ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return s.repo.Create(ctx, req, string(hash))
}

func (s *service) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	rec, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(rec.PasswordHash), []byte(req.OldPassword)); err != nil {
		return ErrInvalidPassword
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, userID, string(newHash))
}

// recordToUser converts an internal DB record to the public-facing User DTO.
func recordToUser(rec *userRecord) *User {
	u := &User{
		ID:          rec.ID,
		Email:       rec.Email,
		NamaLengkap: rec.NamaLengkap,
		NoHp:        rec.NoHp,
		FotoURL:     rec.FotoURL,
		RoleID:      rec.RoleID,
		RoleName:    rec.RoleName,
		Scope:       rec.Scope,
		Status:      rec.Status,
	}
	u.CabangID = rec.CabangID
	u.UnitID = rec.UnitID
	return u
}
