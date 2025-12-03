package handlers

import (
	"encoding/json"
	"net/http"

	pkghttp "github.com/example/back-end-tcc/pkg/http"
	"github.com/example/back-end-tcc/services/trace/service"
)

// HTTP exposes trace endpoints.
type HTTP struct {
	service *service.Service
}

// New creates handlers.
func New(service *service.Service) *HTTP {
	return &HTTP{service: service}
}

// Record stores a new trace event.
func (h *HTTP) Record(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		SubmissionID string `json:"submission_id"`
		Message      string `json:"message"`
		Level        string `json:"level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		pkghttp.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}
	event := h.service.Record(r.Context(), payload.SubmissionID, payload.Message, payload.Level)
	pkghttp.JSON(w, http.StatusCreated, event)
}

// List returns stored events.
func (h *HTTP) List(w http.ResponseWriter, r *http.Request) {
	pkghttp.JSON(w, http.StatusOK, h.service.Events())
}
