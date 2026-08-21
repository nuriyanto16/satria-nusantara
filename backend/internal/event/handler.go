package event

import (
	"github.com/go-chi/chi/v5"
	"satria-nusantara/backend/pkg/response"
	"net/http"
)

type Handler struct{}
func NewHandler() *Handler { return &Handler{} }
func (h *Handler) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			response.Success(w, http.StatusOK, "Event API", []Event{})
		})
		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			response.Success(w, http.StatusOK, "Berhasil mendaftar Latgab", map[string]string{"status": "success"})
		})
		r.Post("/reservasi", func(w http.ResponseWriter, r *http.Request) {
			response.Success(w, http.StatusOK, "Berhasil reservasi", map[string]string{"status": "success"})
		})
	}
}
