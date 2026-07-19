package organization

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"satria-nusantara/backend/pkg/middleware"
	"satria-nusantara/backend/pkg/response"
)

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(jwtSecret string) func(r chi.Router) {
	auth := middleware.Authenticate(jwtSecret)
	pusat := middleware.RequireScope("pusat")
	manager := middleware.RequireScope("pusat", "cabang")

	return func(r chi.Router) {
		r.Use(auth)

		// ── Cabang ─────────────────────────────────────
		r.Get("/cabang", h.listCabang)
		r.Get("/cabang/{id}", h.getCabang)
		r.With(pusat).Post("/cabang", h.createCabang)
		r.With(pusat).Put("/cabang/{id}", h.updateCabang)
		r.Get("/cabang/{id}/unit", h.listUnit)

		// ── Unit ───────────────────────────────────────
		r.With(manager).Post("/unit", h.createUnit)
		r.Get("/unit/{id}", h.getUnit)

		// ── Anggota ────────────────────────────────────
		r.Get("/anggota", h.listAnggota)
		r.Post("/anggota", h.createAnggota)
		r.Get("/anggota/{id}", h.getAnggota)
		r.With(manager).Post("/anggota/{id}/verifikasi", h.verifikasiAnggota)
		r.Post("/anggota/{id}/upload", h.uploadFoto)
		r.Get("/anggota/{id}/stats", h.getAnggotaStats)
		r.Post("/anggota/{id}/kebugaran", h.updateAnggotaKebugaran)

		// ── Pelatih ────────────────────────────────────
		r.Get("/pelatih", h.listPelatih)

		// ── Sebaran ────────────────────────────────────
		r.Get("/sebaran", h.sebaranProvinsi)

		// ── Dashboard ──────────────────────────────────
		r.Get("/dashboard-stats", h.dashboardStats)
	}
}

func (h *Handler) listCabang(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	params := parseListParams(r)
	result, err := h.svc.ListCabang(params, claims.Scope, claims.CabangID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) getCabang(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	c, err := h.svc.GetCabang(id)
	if handleNotFound(w, err) { return }
	response.JSON(w, http.StatusOK, c)
}

func (h *Handler) createCabang(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	var req CreateCabangRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	id, err := h.svc.CreateCabang(req, claims.UserID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, "Cabang berhasil dibuat", map[string]string{"id": id})
}

func (h *Handler) updateCabang(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req CreateCabangRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	if err := h.svc.UpdateCabang(id, req); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, "Cabang berhasil diupdate", nil)
}

func (h *Handler) listUnit(w http.ResponseWriter, r *http.Request) {
	cabangID := chi.URLParam(r, "id")
	units, err := h.svc.ListUnit(cabangID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, units)
}

func (h *Handler) getUnit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	u, err := h.svc.GetUnit(id)
	if handleNotFound(w, err) { return }
	response.JSON(w, http.StatusOK, u)
}

func (h *Handler) createUnit(w http.ResponseWriter, r *http.Request) {
	var req CreateUnitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	id, err := h.svc.CreateUnit(req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, "Unit berhasil dibuat", map[string]string{"id": id})
}

func (h *Handler) listAnggota(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	params.UnitID = r.URL.Query().Get("unit_id")
	result, err := h.svc.ListAnggota(params)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) getAnggota(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	a, err := h.svc.GetAnggota(id)
	if handleNotFound(w, err) { return }
	response.JSON(w, http.StatusOK, a)
}

func (h *Handler) createAnggota(w http.ResponseWriter, r *http.Request) {
	var req CreateAnggotaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	id, err := h.svc.CreateAnggota(req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, "Anggota berhasil didaftarkan, menunggu verifikasi", map[string]string{"id": id})
}

func (h *Handler) verifikasiAnggota(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req VerifikasiAnggotaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	if err := h.svc.VerifikasiAnggota(id, req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	msg := "Anggota berhasil diaktifkan"
	if req.Aksi == "reject" { msg = "Anggota ditolak" }
	response.Success(w, http.StatusOK, msg, nil)
}

func (h *Handler) listPelatih(w http.ResponseWriter, r *http.Request) {
	cabangID := r.URL.Query().Get("cabang_id")
	pelatih, err := h.svc.ListPelatih(cabangID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, pelatih)
}

func (h *Handler) sebaranProvinsi(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.GetSebaranProvinsi()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) dashboardStats(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.GetDashboardStats()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

// ─── Helpers ─────────────────────────────────────────────────────────────────


func parseListParams(r *http.Request) ListParams {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 20 }
	return ListParams{
		Page: page, Limit: limit,
		Search: q.Get("search"), Status: q.Get("status"),
		CabangID: q.Get("cabang_id"),
	}
}

func handleNotFound(w http.ResponseWriter, err error) bool {
	if err == nil { return false }
	if errors.Is(err, ErrNotFound) {
		response.Error(w, http.StatusNotFound, err.Error())
	} else {
		response.Error(w, http.StatusInternalServerError, err.Error())
	}
	return true
}

func (h *Handler) uploadFoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "File terlalu besar atau tidak valid")
		return
	}

	file, header, err := r.FormFile("foto")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "File foto tidak ditemukan")
		return
	}
	defer file.Close()

	// Ensure uploads directory exists inside context
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal membuat direktori upload")
		return
	}

	// Create new file
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("anggota-%s%s", id, ext)
	dstPath := filepath.Join("./uploads", filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal menyimpan file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal menyalin file")
		return
	}

	// Update foto_url in database (absolute URL pointing to backend host)
	fotoURL := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)
	if err := h.svc.UpdateFotoAnggota(id, fotoURL); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Foto berhasil diunggah", map[string]string{"foto_url": fotoURL})
}

func (h *Handler) getAnggotaStats(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	stats, err := h.svc.GetAnggotaStats(id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, stats)
}

// Ensure context is threaded properly — wrapper for service calls
var _ = context.Background


func (h *Handler) updateAnggotaKebugaran(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateKebugaranRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Format payload tidak valid")
		return
	}
	if err := h.svc.UpdateAnggotaKebugaran(id, req); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, "Hasil tes kebugaran berhasil disimpan", nil)
}
