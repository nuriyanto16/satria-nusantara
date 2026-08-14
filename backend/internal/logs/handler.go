package logs

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"satria-nusantara/backend/pkg/response"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/audit", h.getAuditLogs)
		r.Get("/payments", h.getPaymentLogs)
		r.Get("/errors", h.getAppErrors)
	}
}

func (h *Handler) getAuditLogs(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 { page = 1 }
	limit := 20
	offset := (page - 1) * limit

	rows, err := h.db.QueryContext(r.Context(), `
		SELECT a.id, a.user_id, u.nama, a.action, a.entity, a.entity_id, a.details, a.ip_address, a.created_at
		FROM audit_logs a
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY a.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal memuat log audit")
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, action, entity, entityID, ipAddress, userName sql.NullString
		var userID sql.NullString
		var details sql.NullString
		var createdAt time.Time
		if err := rows.Scan(&id, &userID, &userName, &action, &entity, &entityID, &details, &ipAddress, &createdAt); err == nil {
			var d map[string]interface{}
			if details.Valid {
				json.Unmarshal([]byte(details.String), &d)
			}
			result = append(result, map[string]interface{}{
				"id":         id.String,
				"userId":     userID.String,
				"userName":   userName.String,
				"action":     action.String,
				"entity":     entity.String,
				"entityId":   entityID.String,
				"details":    d,
				"ipAddress":  ipAddress.String,
				"createdAt":  createdAt,
			})
		}
	}
	if result == nil { result = []map[string]interface{}{} }
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": result,
		"page": page,
	})
}

func (h *Handler) getPaymentLogs(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 { page = 1 }
	limit := 20
	offset := (page - 1) * limit

	rows, err := h.db.QueryContext(r.Context(), `
		SELECT id, transaction_id, provider, endpoint, request_payload, response_payload, status_code, error_message, created_at
		FROM payment_logs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal memuat log payment")
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, txID, provider, endpoint, errorMsg sql.NullString
		var reqPayload, resPayload sql.NullString
		var statusCode sql.NullInt64
		var createdAt time.Time
		if err := rows.Scan(&id, &txID, &provider, &endpoint, &reqPayload, &resPayload, &statusCode, &errorMsg, &createdAt); err == nil {
			var req, res interface{}
			if reqPayload.Valid { json.Unmarshal([]byte(reqPayload.String), &req) }
			if resPayload.Valid { json.Unmarshal([]byte(resPayload.String), &res) }
			result = append(result, map[string]interface{}{
				"id":               id.String,
				"transactionId":    txID.String,
				"provider":         provider.String,
				"endpoint":         endpoint.String,
				"requestPayload":   req,
				"responsePayload":  res,
				"statusCode":       statusCode.Int64,
				"errorMessage":     errorMsg.String,
				"createdAt":        createdAt,
			})
		}
	}
	if result == nil { result = []map[string]interface{}{} }
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": result,
		"page": page,
	})
}

func (h *Handler) getAppErrors(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 { page = 1 }
	limit := 20
	offset := (page - 1) * limit

	rows, err := h.db.QueryContext(r.Context(), `
		SELECT id, service, error_type, message, stack_trace, context, created_at
		FROM app_errors
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal memuat log error")
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, service, errorType, msg, stack, ctx sql.NullString
		var createdAt time.Time
		if err := rows.Scan(&id, &service, &errorType, &msg, &stack, &ctx, &createdAt); err == nil {
			var c map[string]interface{}
			if ctx.Valid { json.Unmarshal([]byte(ctx.String), &c) }
			result = append(result, map[string]interface{}{
				"id":          id.String,
				"service":     service.String,
				"errorType":   errorType.String,
				"message":     msg.String,
				"stackTrace":  stack.String,
				"context":     c,
				"createdAt":   createdAt,
			})
		}
	}
	if result == nil { result = []map[string]interface{}{} }
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": result,
		"page": page,
	})
}
