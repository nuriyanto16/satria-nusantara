package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"satria-nusantara/backend/internal/logs"
)

// AuditLogMiddleware intercepts requests and logs administrative actions
func AuditLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We only want to log POST, PUT, DELETE for administrative actions
		if r.Method == "GET" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// Read UserID from context if AuthMiddleware was applied before this
		userID := ""
		if val := r.Context().Value("user_id"); val != nil {
			if id, ok := val.(string); ok {
				userID = id
			}
		}

		// For simplicity, entity is derived from the URL path
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		entity := "system"
		if len(pathParts) > 2 {
			entity = pathParts[2] // e.g. /api/v1/admin/users -> users
		} else if len(pathParts) > 1 {
			entity = pathParts[1]
		}

		action := r.Method
		entityID := "" // You could extract ID from URL params if needed

		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // restore body
		}

		// Use the logs service to record the audit log asynchronously
		go func(u, a, e, eid, d, ip string) {
			svc := logs.GetService()
			if svc != nil {
				svc.LogAudit(u, a, e, eid, d, ip)
			}
		}(userID, action, entity, entityID, string(bodyBytes), r.RemoteAddr)

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		// Optionally log application errors if status code is >= 500
		if ww.Status() >= 500 {
			go func(serviceName, errorType, msg, stack, ip string) {
				svc := logs.GetService()
				if svc != nil {
					svc.LogError(serviceName, errorType, msg, stack, map[string]interface{}{
						"ip": ip,
					})
				}
			}("api", "ServerError", http.StatusText(ww.Status()), "", r.RemoteAddr)
		}
	})
}
