package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"satria-nusantara/backend/pkg/response"
)

// contextKey avoids collisions in context values.
type contextKey string

const ClaimsKey contextKey = "claims"

// Claims is the JWT payload stored in every authenticated request context.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	RoleID int    `json:"role_id"`
	// Scope defines what the user can manage: "pusat", "cabang", "unit", "pelatih"
	Scope    string `json:"scope"`
	CabangID string `json:"cabang_id,omitempty"`
	UnitID   string `json:"unit_id,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT token for a user.
func GenerateToken(userID, email, scope, cabangID, unitID, secret string, roleID int) (string, error) {
	claims := Claims{
		UserID:   userID,
		Email:    email,
		RoleID:   roleID,
		Scope:    scope,
		CabangID: cabangID,
		UnitID:   unitID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// Authenticate is a middleware that validates the Bearer JWT token.
func Authenticate(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.Error(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				response.Error(w, http.StatusUnauthorized, "Token tidak valid atau sudah kadaluarsa")
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaims extracts Claims from the request context.
// Returns nil if no claims found (unauthenticated request).
func GetClaims(r *http.Request) *Claims {
	claims, _ := r.Context().Value(ClaimsKey).(*Claims)
	return claims
}

// RequireScope creates a middleware that only allows specific scopes.
// Example: RequireScope("pusat", "cabang") allows both pusat and cabang users.
func RequireScope(scopes ...string) func(http.Handler) http.Handler {
	scopeSet := make(map[string]struct{}, len(scopes))
	for _, s := range scopes {
		scopeSet[s] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r)
			if claims == nil {
				response.Error(w, http.StatusUnauthorized, "Unauthenticated")
				return
			}
			if _, ok := scopeSet[claims.Scope]; !ok {
				response.Error(w, http.StatusForbidden, "Akses tidak diizinkan untuk role ini")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
