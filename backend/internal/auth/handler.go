package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"satria-nusantara/backend/pkg/middleware"
	"satria-nusantara/backend/pkg/response"
)

// Handler holds HTTP handlers for the auth domain.
type Handler struct {
	svc Service
}

// NewHandler creates a new auth handler.
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Routes registers all auth routes onto the given chi router.
// Public routes (no auth) and protected routes are declared here.
func (h *Handler) Routes(jwtSecret string) func(r chi.Router) {
	return func(r chi.Router) {
		// ── Public ─────────────────────────────────────────
		r.Post("/login", h.login)
		r.Post("/google-login", h.googleLogin)
		r.Get("/health", h.health)

		// ── Protected ──────────────────────────────────────
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(jwtSecret))
			r.Get("/me", h.me)
			r.Post("/register", h.register)          // Only logged-in admin can register users
			r.Post("/change-password", h.changePassword)
		})
	}
}

// POST /api/v1/auth/login
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	result, err := h.svc.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidPassword):
			response.Error(w, http.StatusUnauthorized, err.Error())
		case errors.Is(err, ErrUserInactive):
			response.Error(w, http.StatusForbidden, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan server")
		}
		return
	}

	response.Success(w, http.StatusOK, "Login berhasil", result)
}

// POST /api/v1/auth/google-login
func (h *Handler) googleLogin(w http.ResponseWriter, r *http.Request) {
	var req GoogleLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	result, err := h.svc.GoogleLogin(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserInactive):
			response.Error(w, http.StatusForbidden, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(w, http.StatusOK, "Login Google berhasil", result)
}


// GET /api/v1/auth/me
func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.svc.GetMe(r.Context(), claims.UserID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "User tidak ditemukan")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// POST /api/v1/auth/register
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil || claims.Scope != "pusat" {
		response.Error(w, http.StatusForbidden, "Hanya Admin Pusat yang dapat menambah user")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	id, err := h.svc.Register(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailExists):
			response.Error(w, http.StatusConflict, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Gagal membuat user")
		}
		return
	}

	response.Success(w, http.StatusCreated, "User berhasil dibuat", map[string]string{"id": id})
}

// POST /api/v1/auth/change-password
func (h *Handler) changePassword(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	if err := h.svc.ChangePassword(r.Context(), claims.UserID, req); err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			response.Error(w, http.StatusUnauthorized, "Password lama tidak cocok")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Gagal mengubah password")
		return
	}

	response.Success(w, http.StatusOK, "Password berhasil diubah", nil)
}

// GET /api/v1/auth/health
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok", "module": "auth"})
}
