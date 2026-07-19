package auth

import (
	"context"
	"database/sql"
	"errors"
)

// Repository defines the data access contract for the auth domain.
// Implementing this interface makes service logic fully testable.
type Repository interface {
	FindByEmail(ctx context.Context, email string) (*userRecord, error)
	FindByID(ctx context.Context, id string) (*userRecord, error)
	Create(ctx context.Context, req RegisterRequest, passwordHash string) (string, error)
	UpdatePassword(ctx context.Context, userID, newHash string) error
}

// userRecord is the internal DB-level struct (not exposed via API).
type userRecord struct {
	ID           string
	Email        string
	PasswordHash string
	GoogleID     *string
	NamaLengkap  string
	NoHp         string
	FotoURL      string
	RoleID       int
	RoleName     string
	Scope        string
	CabangID     *string
	UnitID       *string
	Status       string
}

// pgRepository is the PostgreSQL implementation of Repository.
type pgRepository struct {
	db *sql.DB
}

// NewRepository creates a new PostgreSQL-backed auth repository.
func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) FindByEmail(ctx context.Context, email string) (*userRecord, error) {
	var u userRecord
	err := r.db.QueryRowContext(ctx, `
		SELECT
			u.id, u.email, u.password_hash, u.google_id,
			u.nama_lengkap, u.no_hp, COALESCE(u.foto_url,''),
			u.role_id, ro.name, ro.scope,
			u.status
		FROM users u
		JOIN roles ro ON ro.id = u.role_id
		WHERE u.email = $1
	`, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.GoogleID,
		&u.NamaLengkap, &u.NoHp, &u.FotoURL,
		&u.RoleID, &u.RoleName, &u.Scope,
		&u.Status,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *pgRepository) FindByID(ctx context.Context, id string) (*userRecord, error) {
	var u userRecord
	err := r.db.QueryRowContext(ctx, `
		SELECT
			u.id, u.email, u.password_hash, u.google_id,
			u.nama_lengkap, u.no_hp, COALESCE(u.foto_url,''),
			u.role_id, ro.name, ro.scope,
			u.status
		FROM users u
		JOIN roles ro ON ro.id = u.role_id
		WHERE u.id = $1
	`, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.GoogleID,
		&u.NamaLengkap, &u.NoHp, &u.FotoURL,
		&u.RoleID, &u.RoleName, &u.Scope,
		&u.Status,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *pgRepository) Create(ctx context.Context, req RegisterRequest, passwordHash string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, nama_lengkap, no_hp, role_id, status)
		VALUES ($1, $2, $3, $4, $5, 'aktif')
		RETURNING id
	`, req.Email, passwordHash, req.NamaLengkap, req.NoHp, req.RoleID).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *pgRepository) UpdatePassword(ctx context.Context, userID, newHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`,
		newHash, userID,
	)
	return err
}
