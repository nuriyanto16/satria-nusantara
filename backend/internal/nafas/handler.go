package nafas

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"satria-nusantara/backend/pkg/response"
)

type HistoryRecord struct {
	ID          string    `json:"id"`
	AnggotaID   string    `json:"anggota_id"`
	AnggotaNama string    `json:"anggota_nama"`
	CabangNama  string    `json:"cabang_nama"`
	UnitNama    string    `json:"unit_nama"`
	Teknik      string    `json:"teknik"`       // "Sama Kaki", "Sama Sisi", "Segi Empat"
	DurasiDetik int       `json:"durasi_detik"` // e.g. 600
	DurasiFmt   string    `json:"durasi_fmt"`   // e.g. "10:00"
	Siklus      int       `json:"siklus"`       // e.g. 12
	Timestamp   time.Time `json:"timestamp"`
}

type Handler struct {
	mu        sync.RWMutex
	histories []HistoryRecord
}

func NewHandler() *Handler {
	// Seed sample history data for monitoring
	now := time.Now()
	sampleHistories := []HistoryRecord{
		{
			ID:          "hist-101",
			AnggotaID:   "a0000001-0000-0000-0000-000000000001",
			AnggotaNama: "Ahmad Santoso",
			CabangNama:  "Kota Yogyakarta",
			UnitNama:    "Malioboro",
			Teknik:      "Sama Kaki",
			DurasiDetik: 600,
			DurasiFmt:   "10:00",
			Siklus:      12,
			Timestamp:   now.Add(-25 * time.Minute),
		},
		{
			ID:          "hist-102",
			AnggotaID:   "a0000002-0000-0000-0000-000000000002",
			AnggotaNama: "Budi Susanto",
			CabangNama:  "Surabaya",
			UnitNama:    "Rungkut",
			Teknik:      "Segi Empat",
			DurasiDetik: 480,
			DurasiFmt:   "08:00",
			Siklus:      12,
			Timestamp:   now.Add(-2 * time.Hour),
		},
		{
			ID:          "hist-103",
			AnggotaID:   "a0000003-0000-0000-0000-000000000003",
			AnggotaNama: "Siti Rahmawati",
			CabangNama:  "Kota Bandung",
			UnitNama:    "Dago",
			Teknik:      "Sama Sisi",
			DurasiDetik: 300,
			DurasiFmt:   "05:00",
			Siklus:      10,
			Timestamp:   now.Add(-5 * time.Hour),
		},
		{
			ID:          "hist-104",
			AnggotaID:   "a0000004-0000-0000-0000-000000000004",
			AnggotaNama: "Dewi Lestari",
			CabangNama:  "DKI Jakarta",
			UnitNama:    "Monas",
			Teknik:      "Sama Kaki",
			DurasiDetik: 900,
			DurasiFmt:   "15:00",
			Siklus:      18,
			Timestamp:   now.Add(-24 * time.Hour),
		},
		{
			ID:          "hist-105",
			AnggotaID:   "a0000001-0000-0000-0000-000000000001",
			AnggotaNama: "Ahmad Santoso",
			CabangNama:  "Kota Yogyakarta",
			UnitNama:    "Malioboro",
			Teknik:      "Sama Kaki",
			DurasiDetik: 600,
			DurasiFmt:   "10:00",
			Siklus:      12,
			Timestamp:   now.Add(-48 * time.Hour),
		},
	}

	return &Handler{
		histories: sampleHistories,
	}
}

func (h *Handler) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/history", h.saveHistory)
		r.Get("/history", h.getHistory)
		r.Get("/monitoring", h.getMonitoringData)
	}
}

func (h *Handler) saveHistory(w http.ResponseWriter, r *http.Request) {
	var req HistoryRecord
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if req.ID == "" {
		req.ID = "hist-" + time.Now().Format("20060102150405")
	}
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}
	if req.AnggotaNama == "" {
		req.AnggotaNama = "Anggota Satria Nusantara"
	}
	if req.CabangNama == "" {
		req.CabangNama = "Pusat"
	}
	if req.DurasiFmt == "" && req.DurasiDetik > 0 {
		m := req.DurasiDetik / 60
		s := req.DurasiDetik % 60
		req.DurasiFmt = time.Date(0, 0, 0, 0, m, s, 0, time.UTC).Format("04:05")
	}

	h.histories = append([]HistoryRecord{req}, h.histories...)

	response.Success(w, http.StatusCreated, "Riwayat latihan berhasil disimpan", req)
}

func (h *Handler) getHistory(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	anggotaID := r.URL.Query().Get("anggota_id")
	search := strings.ToLower(r.URL.Query().Get("search"))

	var filtered []HistoryRecord
	for _, rec := range h.histories {
		if anggotaID != "" && rec.AnggotaID != anggotaID {
			continue
		}
		if search != "" {
			nameMatch := strings.Contains(strings.ToLower(rec.AnggotaNama), search)
			teknikMatch := strings.Contains(strings.ToLower(rec.Teknik), search)
			if !nameMatch && !teknikMatch {
				continue
			}
		}
		filtered = append(filtered, rec)
	}

	response.JSON(w, http.StatusOK, filtered)
}

func (h *Handler) getMonitoringData(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	search := strings.ToLower(r.URL.Query().Get("search"))
	teknikFilter := strings.ToLower(r.URL.Query().Get("teknik"))

	var filtered []HistoryRecord
	totalSec := 0
	uniqueAnggota := make(map[string]bool)
	teknikCount := make(map[string]int)

	for _, rec := range h.histories {
		if teknikFilter != "" && strings.ToLower(rec.Teknik) != teknikFilter {
			continue
		}
		if search != "" {
			nameMatch := strings.Contains(strings.ToLower(rec.AnggotaNama), search)
			cabangMatch := strings.Contains(strings.ToLower(rec.CabangNama), search)
			if !nameMatch && !cabangMatch {
				continue
			}
		}

		filtered = append(filtered, rec)
		totalSec += rec.DurasiDetik
		uniqueAnggota[rec.AnggotaID] = true
		teknikCount[rec.Teknik]++
	}

	favTeknik := "Sama Kaki"
	maxCount := 0
	for tek, cnt := range teknikCount {
		if cnt > maxCount {
			maxCount = cnt
			favTeknik = tek
		}
	}

	stats := map[string]interface{}{
		"total_sesi":       len(filtered),
		"total_menit":      totalSec / 60,
		"anggota_aktif":    len(uniqueAnggota),
		"teknik_favorit":   favTeknik,
		"histories":        filtered,
	}

	response.JSON(w, http.StatusOK, stats)
}
