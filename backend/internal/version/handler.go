package version

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"satria-nusantara/backend/pkg/response"
)

type AppVersion struct {
	ID          string    `json:"id"`
	VersionName string    `json:"version_name"`
	BuildNumber int       `json:"build_number"`
	ReleaseNotes string   `json:"release_notes"`
	IsMandatory bool      `json:"is_mandatory"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateVersionRequest struct {
	VersionName string `json:"version_name"`
	BuildNumber int    `json:"build_number"`
	ReleaseNotes string `json:"release_notes"`
	IsMandatory bool   `json:"is_mandatory"`
}

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RoutesAdmin(r chi.Router) {
	// Admin routes
	r.Get("/versions", h.listVersions)
	r.Post("/versions", h.createVersion)
	r.Delete("/versions/{id}", h.deleteVersion)

	// Public route
	
}

func (h *Handler) listVersions(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(r.Context(), "SELECT id, version_name, build_number, release_notes, is_mandatory, created_at FROM app_versions ORDER BY build_number DESC")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var res []AppVersion
	for rows.Next() {
		var v AppVersion
		if err := rows.Scan(&v.ID, &v.VersionName, &v.BuildNumber, &v.ReleaseNotes, &v.IsMandatory, &v.CreatedAt); err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		res = append(res, v)
	}
	response.Success(w, http.StatusOK, "Berhasil memuat daftar versi", res)
}

func (h *Handler) createVersion(w http.ResponseWriter, r *http.Request) {
	var req CreateVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	
	var id string
	err := h.db.QueryRowContext(r.Context(), "INSERT INTO app_versions (version_name, build_number, release_notes, is_mandatory) VALUES ($1, $2, $3, $4) RETURNING id", req.VersionName, req.BuildNumber, req.ReleaseNotes, req.IsMandatory).Scan(&id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, "Versi berhasil ditambahkan", map[string]string{"id": id})
}

func (h *Handler) deleteVersion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.ExecContext(r.Context(), "DELETE FROM app_versions WHERE id = $1", id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, "Versi berhasil dihapus", nil)
}

func (h *Handler) CheckVersion(w http.ResponseWriter, r *http.Request) {
	var v AppVersion
	err := h.db.QueryRowContext(r.Context(), "SELECT id, version_name, build_number, release_notes, is_mandatory, created_at FROM app_versions ORDER BY build_number DESC LIMIT 1").Scan(&v.ID, &v.VersionName, &v.BuildNumber, &v.ReleaseNotes, &v.IsMandatory, &v.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Success(w, http.StatusOK, "Belum ada versi dirilis", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, "Versi terbaru", v)
}
