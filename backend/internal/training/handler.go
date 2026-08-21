package training

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"satria-nusantara/backend/pkg/middleware"
	"satria-nusantara/backend/pkg/response"
)

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(jwtSecret string) func(r chi.Router) {
	auth := middleware.Authenticate(jwtSecret)

	return func(r chi.Router) {
		r.Use(auth)
		
		r.Get("/sesi", h.listSesi)
		r.Post("/sesi", h.createSesi)
		r.Put("/sesi/{id}", h.updateSesi)
		r.Delete("/sesi/{id}", h.deleteSesi)
		
		r.Post("/sesi/{id}/qr", h.generateQR)
		r.Post("/scan-qr", h.scanQR)
		
		r.Get("/antrian/status", h.getAntrianStatus)
		
		r.Get("/sesi/{id}/kehadiran", h.getKehadiran)
	}
}

func (h *Handler) listSesi(w http.ResponseWriter, r *http.Request) {
	unitID := r.URL.Query().Get("unit_id")
	tanggal := r.URL.Query().Get("tanggal")
	sesi, err := h.svc.ListSesi(unitID, tanggal)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, sesi)
}

func (h *Handler) createSesi(w http.ResponseWriter, r *http.Request) {
	var req CreateSesiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	id, err := h.svc.CreateSesi(req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, "Sesi berhasil dibuat", map[string]string{"id": id})
}

func (h *Handler) deleteSesi(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteSesi(id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, "Sesi berhasil dihapus", nil)
}

func (h *Handler) updateSesi(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req CreateSesiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.svc.UpdateSesi(id, req); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, "Sesi berhasil diperbarui", nil)
}

func (h *Handler) generateQR(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	qr, err := h.svc.GenerateQR(id, 4) // default 4 hours
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, qr)
}

func (h *Handler) scanQR(w http.ResponseWriter, r *http.Request) {
	var req ScanQRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.svc.ScanQR(req); err != nil {
		if err == ErrQRExpired || err == ErrAlreadyScan || err == ErrNotFound {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, "Berhasil absen", nil)
}

func (h *Handler) getKehadiran(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ringkasan, detail, err := h.svc.GetRingkasan(id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{
		"ringkasan": ringkasan,
		"detail": detail,
	})
}

func (h *Handler) getAntrianStatus(w http.ResponseWriter, r *http.Request) {
	// Dummy response for antrian status
	data := map[string]any{
		"nomor_antrian": 14,
		"antrian_saat_ini": 9,
		"total_antrian": 25,
		"estimasi": "± 25 menit",
		"nama_session": "Latihan Gabungan – Minggu, 10 Agu 2026",
		"lokasi": "GOR Arcamanik, Kota Bandung",
		"antrian": []map[string]any{
			{"no": 1, "nama": "Ahmad Fauzi", "status": "selesai"},
			{"no": 2, "nama": "Budi Santoso", "status": "selesai"},
			{"no": 3, "nama": "Citra Dewi", "status": "selesai"},
			{"no": 4, "nama": "Dani Kurniawan", "status": "selesai"},
			{"no": 5, "nama": "Eva Susanti", "status": "selesai"},
			{"no": 6, "nama": "Fajar Nugroho", "status": "selesai"},
			{"no": 7, "nama": "Galih Prasetyo", "status": "selesai"},
			{"no": 8, "nama": "Hani Pertiwi", "status": "selesai"},
			{"no": 9, "nama": "Irfan Hakim", "status": "aktif"},
			{"no": 10, "nama": "Joko Susilo", "status": "menunggu"},
			{"no": 11, "nama": "Kartika Sari", "status": "menunggu"},
			{"no": 12, "nama": "Lukman Hakim", "status": "menunggu"},
			{"no": 13, "nama": "Mira Lesmana", "status": "menunggu"},
			{"no": 14, "nama": "Anda", "status": "menunggu"},
			{"no": 15, "nama": "Nina Bobo", "status": "menunggu"},
		},
	}
	response.Success(w, http.StatusOK, "Berhasil mengambil status antrian", data)
}
