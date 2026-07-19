// Package response provides a consistent JSON response format for all API endpoints.
//
// All responses follow the structure:
//
//	{ "data": ..., "message": "...", "error": "..." }
package response

import (
	"encoding/json"
	"net/http"
)

// envelope is the standard JSON response wrapper.
type envelope struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// JSON writes a successful JSON response.
func JSON(w http.ResponseWriter, status int, data any) {
	write(w, status, envelope{Data: data})
}

// Success writes a success response with a message and optional data.
func Success(w http.ResponseWriter, status int, message string, data any) {
	write(w, status, envelope{Data: data, Message: message})
}

// Error writes an error response.
func Error(w http.ResponseWriter, status int, message string) {
	write(w, status, envelope{Error: message})
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func write(w http.ResponseWriter, status int, body envelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
