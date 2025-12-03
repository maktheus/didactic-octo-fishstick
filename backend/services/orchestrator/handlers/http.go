package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	pkghttp "github.com/example/back-end-tcc/pkg/http"
	"github.com/example/back-end-tcc/services/orchestrator/service"
)

// HTTP provides orchestrator endpoints.
type HTTP struct {
	service *service.Service
}

// New creates handlers.
func New(service *service.Service) *HTTP {
	return &HTTP{service: service}
}

// Submit enqueues a submission.
func (h *HTTP) Submit(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		BenchmarkID string `json:"benchmark_id"`
		AgentID     string `json:"agent_id"`
		Payload     string `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		pkghttp.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}
	submission, err := h.service.Submit(context.Background(), payload.BenchmarkID, payload.AgentID, payload.Payload)
	if err != nil {
		pkghttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	pkghttp.JSON(w, http.StatusAccepted, submission)
}

// List returns submissions.
func (h *HTTP) List(w http.ResponseWriter, r *http.Request) {
	pkghttp.JSON(w, http.StatusOK, h.service.List())
}
