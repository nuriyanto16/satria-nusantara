package finance

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"satria-nusantara/backend/pkg/response"
)

// Transaction represents a payment transaction
type Transaction struct {
	ID            string    `json:"id"`
	IuranID       string    `json:"iuranId"`
	UserID        string    `json:"userId"`
	Nama          string    `json:"nama"`
	Nomor         string    `json:"nomor"`
	Bulan         string    `json:"bulan"`
	Amount        int       `json:"amount"`
	PaymentMethod string    `json:"paymentMethod"`
	Status        string    `json:"status"` // 'pending', 'lunas', 'ditolak'
	BuktiUrl      string    `json:"buktiUrl"`
	CreatedAt     time.Time `json:"createdAt"`
}

// In-memory data store for transactions
var (
	transactions = []Transaction{
		{
			ID:            "pay_1753145520123",
			IuranID:       "2",
			UserID:        "u_sari",
			Nama:          "Sari Rahmawati",
			Nomor:         "YO-YGY-00098",
			Bulan:         "Juli 2026",
			Amount:        40000,
			PaymentMethod: "Transfer Bank - BCA",
			Status:        "pending",
			BuktiUrl:      "bukti_transfer_sample.jpg",
			CreatedAt:     time.Now().Add(-5 * time.Hour),
		},
		{
			ID:            "pay_1753098610234",
			IuranID:       "4",
			UserID:        "u_agus",
			Nama:          "Agus Prasetyo",
			Nomor:         "YO-YGY-00201",
			Bulan:         "Juli 2026",
			Amount:        40000,
			PaymentMethod: "Transfer Bank - BRI",
			Status:        "pending",
			BuktiUrl:      "bukti_transfer_sample.jpg",
			CreatedAt:     time.Now().Add(-19 * time.Hour),
		},
		{
			ID:            "pay_1752925320456",
			IuranID:       "7",
			UserID:        "u_nurul",
			Nama:          "Nurul Rahayu",
			Nomor:         "YO-YGY-00089",
			Bulan:         "Juli 2026",
			Amount:        40000,
			PaymentMethod: "Transfer Bank - Mandiri",
			Status:        "ditolak",
			BuktiUrl:      "bukti_transfer_sample.jpg",
			CreatedAt:     time.Now().Add(-4 * 24 * time.Hour),
		},
	}
	mu sync.Mutex
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/iuran", h.getIuranHistory)
		r.Post("/iuran/pay", h.payIuran)
	}
}

func (h *Handler) getIuranHistory(w http.ResponseWriter, r *http.Request) {
	// Simple stub response for iuran
	response.Success(w, http.StatusOK, "Finance BLBA API", []BLBA{})
}

type PayRequest struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	UserID string `json:"userId"`
}

func (h *Handler) payIuran(w http.ResponseWriter, r *http.Request) {
	var req PayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	status := "lunas"
	var buktiUrl string
	if strings.HasPrefix(req.Method, "Transfer") {
		status = "pending"
		buktiUrl = "bukti_transfer_sample.jpg" // mock uploaded proof file name
	}

	// Try to infer name and member number based on userId
	nama := "Anggota SN"
	nomor := "YO-YGY-00142"
	if req.UserID != "" {
		if strings.Contains(req.UserID, "registered") {
			nama = "Pendaftar Baru"
			nomor = "SN-NEWUSER"
		} else if strings.Contains(req.UserID, "google") {
			nama = "User Google"
			nomor = "SN-GOOGLE"
		}
	}

	newTrx := Transaction{
		ID:            "pay_" + time.Now().Format("20060102150405"),
		IuranID:       req.ID,
		UserID:        req.UserID,
		Nama:          nama,
		Nomor:         nomor,
		Bulan:         "Juli 2026",
		Amount:        40000,
		PaymentMethod: req.Method,
		Status:        status,
		BuktiUrl:      buktiUrl,
		CreatedAt:     time.Now(),
	}

	transactions = append([]Transaction{newTrx}, transactions...)

	response.Success(w, http.StatusOK, "Payment recorded successfully", newTrx)
}

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(transactions)
}

func (h *Handler) VerifyTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	var req struct {
		Status string `json:"status"` // 'lunas' or 'ditolak'
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	found := false
	for i, t := range transactions {
		if t.ID == id {
			transactions[i].Status = req.Status
			found = true
			break
		}
	}

	if !found {
		response.Error(w, http.StatusNotFound, "Transaction not found")
		return
	}

	response.Success(w, http.StatusOK, "Transaction status updated successfully", nil)
}
