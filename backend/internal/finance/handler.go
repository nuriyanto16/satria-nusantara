package finance

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"satria-nusantara/backend/pkg/response"
)

// PaymentTransaction represents a persisted payment transaction
type PaymentTransaction struct {
	ID            string     `json:"id"`
	ReferenceType string     `json:"referenceType"`
	ReferenceID   string     `json:"referenceId"`
	UserID        string     `json:"userId"`
	Nama          string     `json:"nama"`
	Nomor         string     `json:"nomor"`
	Amount        int        `json:"amount"`
	Description   string     `json:"description"`
	Provider      string     `json:"provider"`
	ProviderID    string     `json:"providerId"`
	PaymentURL    string     `json:"paymentUrl"`
	Status        string     `json:"status"`
	PaidAt        *time.Time `json:"paidAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
}

type Handler struct {
	db     *sql.DB
	xendit *XenditService
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db:     db,
		xendit: NewXenditService(),
	}
}

func (h *Handler) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/history", h.getPaymentHistory) // For mobile
		r.Get("/iuran", h.getIuranHistory)
		r.Post("/iuran/pay", h.payIuran)
		r.Post("/iuran/xendit/create-invoice", h.createXenditInvoice)
		
		r.Post("/webhook/xendit/invoice", h.webhookXenditInvoice)
		r.Get("/webhook/xendit/invoice", h.webhookXenditTest)
		r.Post("/webhook/xendit/ewallet", h.webhookXenditEwallet)
		r.Get("/webhook/xendit/ewallet", h.webhookXenditTest)
		r.Post("/webhook/xendit/payment-request", h.webhookXenditPaymentRequest)
		r.Get("/webhook/xendit/payment-request", h.webhookXenditTest)
	}
}

// GET /api/v1/admin/finance/transactions
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		response.Success(w, http.StatusOK, "Transactions API (No DB)", []PaymentTransaction{})
		return
	}

	rows, err := h.db.QueryContext(r.Context(), `
		SELECT p.id, p.reference_type, p.reference_id, p.user_id, p.amount, p.description, p.provider, p.provider_id, p.payment_url, p.status, p.paid_at, p.created_at, u.nama_lengkap, u.nomor_anggota
		FROM payment_transactions p
		LEFT JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to query transactions")
		return
	}
	defer rows.Close()

	var txs []PaymentTransaction
	for rows.Next() {
		var tx PaymentTransaction
		var refID, uID, nama, nomor sql.NullString
		var paidAt sql.NullTime
		if err := rows.Scan(&tx.ID, &tx.ReferenceType, &refID, &uID, &tx.Amount, &tx.Description, &tx.Provider, &tx.ProviderID, &tx.PaymentURL, &tx.Status, &paidAt, &tx.CreatedAt, &nama, &nomor); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		if refID.Valid { tx.ReferenceID = refID.String }
		if uID.Valid { tx.UserID = uID.String }
		if nama.Valid { tx.Nama = nama.String }
		if nomor.Valid { tx.Nomor = nomor.String }
		if paidAt.Valid { tx.PaidAt = &paidAt.Time }
		txs = append(txs, tx)
	}
	if txs == nil { txs = []PaymentTransaction{} }

	response.Success(w, http.StatusOK, "Admin Transactions", txs)
}

// POST /api/v1/admin/iuran-transactions/{id}/sync
func (h *Handler) SyncXenditTransaction(w http.ResponseWriter, r *http.Request) {
	txID := chi.URLParam(r, "id")
	if h.db == nil {
		response.Error(w, http.StatusInternalServerError, "No DB connection")
		return
	}

	var providerID string
	err := h.db.QueryRow(`SELECT provider_id FROM payment_transactions WHERE id = $1`, txID).Scan(&providerID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Transaction not found")
		return
	}

	// Fetch from Xendit
	inv, err := h.xendit.GetInvoice(providerID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal sinkronisasi dari Xendit: " + err.Error())
		return
	}

	// Update DB
	if inv.Status == "PAID" {
		var paidAt time.Time
		if inv.PaidAt != nil {
			paidAt, _ = time.Parse(time.RFC3339, *inv.PaidAt)
		} else {
			paidAt = time.Now()
		}
		h.db.Exec(`UPDATE payment_transactions SET status = 'PAID', paid_at = $1 WHERE id = $2`, paidAt, txID)
	} else if inv.Status == "EXPIRED" {
		h.db.Exec(`UPDATE payment_transactions SET status = 'EXPIRED' WHERE id = $1`, txID)
	}

	response.Success(w, http.StatusOK, "Sinkronisasi berhasil", inv)
}

// GET /api/v1/finance/history?userId=xxx
func (h *Handler) getPaymentHistory(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if h.db == nil || userId == "" {
		response.Success(w, http.StatusOK, "Payment History", []PaymentTransaction{})
		return
	}

	rows, err := h.db.QueryContext(r.Context(), `
		SELECT id, reference_type, reference_id, user_id, amount, description, provider, provider_id, payment_url, status, paid_at, created_at
		FROM payment_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to query history")
		return
	}
	defer rows.Close()

	var txs []PaymentTransaction
	for rows.Next() {
		var tx PaymentTransaction
		var refID sql.NullString
		var uID sql.NullString
		var paidAt sql.NullTime
		if err := rows.Scan(&tx.ID, &tx.ReferenceType, &refID, &uID, &tx.Amount, &tx.Description, &tx.Provider, &tx.ProviderID, &tx.PaymentURL, &tx.Status, &paidAt, &tx.CreatedAt); err != nil {
			continue
		}
		if refID.Valid { tx.ReferenceID = refID.String }
		if uID.Valid { tx.UserID = uID.String }
		if paidAt.Valid { tx.PaidAt = &paidAt.Time }
		txs = append(txs, tx)
	}
	if txs == nil { txs = []PaymentTransaction{} }

	response.Success(w, http.StatusOK, "Payment History", txs)
}

func (h *Handler) webhookXenditTest(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "Webhook endpoint is active. Please use POST.", nil)
}

func (h *Handler) getIuranHistory(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "Finance BLBA API", []interface{}{})
}

type PayRequest struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	UserID string `json:"userId"`
	Bulan  string `json:"bulan,omitempty"`
	Amount int    `json:"amount,omitempty"`
}

var mu sync.Mutex // For legacy payIuran stub

func (h *Handler) payIuran(w http.ResponseWriter, r *http.Request) {
	var req PayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	response.Success(w, http.StatusOK, "Pembayaran berhasil dicatat (Stub)", nil)
}

type CreateXenditInvoiceReq struct {
	TransactionID string `json:"transactionId"` 
	UserID        string `json:"userId"`
	Nama          string `json:"nama"`
	Email         string `json:"email"`
	Bulan         string `json:"bulan"`
	Amount        int    `json:"amount"`
}

func (h *Handler) createXenditInvoice(w http.ResponseWriter, r *http.Request) {
	var req CreateXenditInvoiceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Amount <= 0 { req.Amount = 40000 }

	externalID := fmt.Sprintf("SN-BLBA-%s-%d", req.UserID, time.Now().UnixMilli())
	desc := fmt.Sprintf("BLBA Satria Nusantara - %s", req.Bulan)

	invoice, err := h.xendit.CreateInvoice(CreateInvoiceRequest{
		ExternalID:         externalID,
		Amount:             req.Amount,
		Description:        desc,
		PayerEmail:         req.Email,
		Currency:           "IDR",
		InvoiceDuration:    86400,
		SuccessRedirectURL: "satrianusantara://payment-success",
		FailureRedirectURL: "satrianusantara://payment-success",
		PaymentMethods:     []string{"OVO", "DANA", "LINKAJA", "SHOPEEPAY", "BCA", "BNI", "BRI", "MANDIRI", "PERMATA", "ALFAMART", "INDOMARET"},
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal membuat invoice Xendit: "+err.Error())
		return
	}

	if h.db != nil {
		// handle non-uuid gracefully if testing
		userIDVal := sql.NullString{String: req.UserID, Valid: true}
		if len(req.UserID) != 36 { userIDVal.Valid = false } // simple UUID length check for mock users

		_, err := h.db.ExecContext(r.Context(), `
			INSERT INTO payment_transactions (reference_type, user_id, amount, description, provider, provider_id, payment_url, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, "blba", userIDVal, req.Amount, desc, "xendit", invoice.ID, invoice.InvoiceURL, "PENDING")
		if err != nil {
			log.Println("Failed to insert payment transaction:", err)
		}
	}

	response.Success(w, http.StatusOK, "Invoice berhasil dibuat", map[string]interface{}{
		"invoiceUrl": invoice.InvoiceURL,
		"invoiceId":  invoice.ID,
		"externalId": invoice.ExternalID,
	})
}

func (h *Handler) VerifyTransaction(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "Transaksi telah diverifikasi (Stub)", nil)
}

func (h *Handler) webhookXenditInvoice(w http.ResponseWriter, r *http.Request) {
	if !h.xendit.VerifyWebhookToken(r) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"invalid webhook token"}`))
		return
	}

	body, _ := io.ReadAll(r.Body)
	var payload XenditWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("[Xendit Webhook] Invoice ID: %s, ExternalID: %s, Status: %s", payload.ID, payload.ExternalID, payload.Status)

	if payload.Status == "PAID" && h.db != nil {
		paidAt, _ := time.Parse(time.RFC3339, payload.PaidAt)
		_, err := h.db.Exec(`UPDATE payment_transactions SET status = 'PAID', paid_at = $1 WHERE provider_id = $2`, paidAt, payload.ID)
		if err != nil {
			log.Println("Failed to update payment_transactions:", err)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handler) webhookXenditEwallet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handler) webhookXenditPaymentRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
