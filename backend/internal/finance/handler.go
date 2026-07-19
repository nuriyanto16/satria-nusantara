package finance

import (
	"github.com/go-chi/chi/v5"
	"satria-nusantara/backend/pkg/response"
	"net/http"
)

type Handler struct{}
func NewHandler() *Handler { return &Handler{} }
func (h *Handler) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/blba", func(w http.ResponseWriter, r *http.Request) {
			response.Success(w, http.StatusOK, "Finance BLBA API", []BLBA{})
		})
	}
}
