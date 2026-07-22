// Package auth defines the data models and DTOs for the authentication domain.
package auth

import "time"

// User represents a system user (pengurus/admin).
type User struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	NamaLengkap string     `json:"nama_lengkap"`
	NoHp        string     `json:"no_hp"`
	FotoURL     string     `json:"foto_url"`
	RoleID      int        `json:"role_id"`
	RoleName    string     `json:"role_name"`
	Scope       string     `json:"scope"`
	CabangID    *string    `json:"cabang_id,omitempty"`
	UnitID      *string    `json:"unit_id,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}

// --- Request DTOs ---

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GoogleLoginRequest struct {
	IDToken     string `json:"id_token,omitempty"`
	Email       string `json:"email"`
	GoogleID    string `json:"google_id,omitempty"`
	NamaLengkap string `json:"nama_lengkap,omitempty"`
	FotoURL     string `json:"foto_url,omitempty"`
	NoHp        string `json:"no_hp,omitempty"`
	CabangID    string `json:"cabang_id,omitempty"`
	UnitID      string `json:"unit_id,omitempty"`
}

type SignupAnggotaRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	NamaLengkap string `json:"nama_lengkap"`
	NoHp        string `json:"no_hp"`
	UnitID      string `json:"unit_id"`
	Tingkatan   string `json:"tingkatan"`
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	NamaLengkap string `json:"nama_lengkap"`
	NoHp        string `json:"no_hp"`
	RoleID      int    `json:"role_id"`
	CabangID    string `json:"cabang_id"`
	UnitID      string `json:"unit_id"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// --- Response DTOs ---

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
