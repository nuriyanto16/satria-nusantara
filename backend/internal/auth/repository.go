package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

// Repository defines the data access contract for the auth domain.
// Implementing this interface makes service logic fully testable.
type Repository interface {
	FindByEmail(ctx context.Context, email string) (*userRecord, error)
	FindByID(ctx context.Context, id string) (*userRecord, error)
	Create(ctx context.Context, req RegisterRequest, passwordHash string) (string, error)
	CreatePendingAnggota(ctx context.Context, req SignupAnggotaRequest, passwordHash string) (string, error)
	CreateGoogleUser(ctx context.Context, req GoogleLoginRequest) (*userRecord, error)
	UpdatePassword(ctx context.Context, userID, newHash string) error
	UpdateGoogleID(ctx context.Context, userID, googleID string) error
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

func (r *pgRepository) CreatePendingAnggota(ctx context.Context, req SignupAnggotaRequest, passwordHash string) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// 1. Resolve Unit Latihan ID (by name ILIKE first, then uuid)
	var resolvedUnitID string
	unitSearch := req.UnitID
	if strings.Contains(unitSearch, " · ") {
		unitSearch = strings.Split(unitSearch, " · ")[0]
	}
	unitSearch = strings.TrimSpace(unitSearch)

	// Try searching by name containing unitSearch
	err = tx.QueryRowContext(ctx, "SELECT id FROM unit_latihan WHERE nama ILIKE $1", "%"+unitSearch+"%").Scan(&resolvedUnitID)
	if err != nil {
		// Fallback to exact UUID
		err = tx.QueryRowContext(ctx, "SELECT id FROM unit_latihan WHERE id::text = $1", req.UnitID).Scan(&resolvedUnitID)
		if err != nil {
			// Fallback: cari unit latihan pertama
			err = tx.QueryRowContext(ctx, "SELECT id FROM unit_latihan LIMIT 1").Scan(&resolvedUnitID)
			if err != nil {
				resolvedUnitID = "00000000-0000-0000-0000-000000000000"
			}
		}
	}

	// 2. Resolve Tingkatan ENUM
	tingkatan := "Pra Dasar"
	cleanTingkat := strings.ToLower(req.Tingkatan)
	if strings.Contains(cleanTingkat, "pra dasar") {
		tingkatan = "Pra Dasar"
	} else if strings.Contains(cleanTingkat, "dasar") {
		tingkatan = "Dasar"
	} else if strings.Contains(cleanTingkat, "gpk") {
		tingkatan = "GPK"
	} else if strings.Contains(cleanTingkat, "pk") {
		tingkatan = "PK"
	} else if strings.Contains(cleanTingkat, "ph") {
		tingkatan = "PH"
	} else if strings.Contains(cleanTingkat, "gabungan") {
		tingkatan = "Gabungan"
	} else if strings.Contains(cleanTingkat, "penjuru") {
		tingkatan = "Penjuru"
	} else if strings.Contains(cleanTingkat, "selangkah") {
		tingkatan = "Selangkah"
	} else if strings.Contains(cleanTingkat, "meditasi") {
		tingkatan = "Meditasi"
	}

	// 3. Insert ke tabel users dengan status 'pending'
	var userID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, nama_lengkap, no_hp, role_id, status)
		VALUES ($1, $2, $3, $4, 4, 'pending')
		RETURNING id
	`, req.Email, passwordHash, req.NamaLengkap, req.NoHp).Scan(&userID)
	if err != nil {
		return "", err
	}

	// 4. Parse Tanggal Lahir
	var tglLahir *time.Time
	if req.TanggalLahir != "" {
		t, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err == nil {
			tglLahir = &t
		}
	}

	// 5. Insert ke tabel anggota dengan status 'pending'
	_, err = tx.ExecContext(ctx, `
		INSERT INTO anggota (user_id, nama_lengkap, no_hp, unit_id, tingkatan, tanggal_lahir, jenis_kelamin, status)
		VALUES ($1, $2, $3, $4, $5::tingkatan_enum, $6, $7, 'pending')
	`, userID, req.NamaLengkap, req.NoHp, resolvedUnitID, tingkatan, tglLahir, req.JenisKelamin)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (r *pgRepository) CreateGoogleUser(ctx context.Context, req GoogleLoginRequest) (*userRecord, error) {
	var id string
	googleID := req.GoogleID
	if googleID == "" {
		googleID = "goog_" + req.Email
	}
	nama := req.NamaLengkap
	if nama == "" {
		nama = req.Email
	}
	foto := req.FotoURL
	if foto == "" {
		foto = "https://ui-avatars.com/api/?name=" + nama
	}
	noHp := req.NoHp
	if noHp == "" {
		noHp = "081234567890"
	}

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (email, google_id, nama_lengkap, no_hp, foto_url, role_id, status)
		VALUES ($1, $2, $3, $4, $5, 4, 'aktif')
		RETURNING id
	`, req.Email, googleID, nama, noHp, foto).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *pgRepository) UpdatePassword(ctx context.Context, userID, newHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`,
		newHash, userID,
	)
	return err
}

func (r *pgRepository) UpdateGoogleID(ctx context.Context, userID, googleID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET google_id = $1, updated_at = NOW() WHERE id = $2`,
		googleID, userID,
	)
	return err
}

