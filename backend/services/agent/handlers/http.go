package handlers

import (
	"encoding/json"
	"net/http"

	pkghttp "github.com/example/back-end-tcc/pkg/http"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/services/agent/service"
)

// HTTP exposes agent endpoints.
type HTTP struct {
	service *service.AgentService
}

// NewHTTP creates handlers.
func NewHTTP(service *service.AgentService) *HTTP {
	return &HTTP{service: service}
}

// Register registers an agent via POST body.
func (h *HTTP) Register(w http.ResponseWriter, r *http.Request) {
	var agent models.User
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		pkghttp.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := h.service.Register(&agent); err != nil {
		pkghttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	pkghttp.JSON(w, http.StatusCreated, agent)
}

// List returns the registered agents.
func (h *HTTP) List(w http.ResponseWriter, r *http.Request) {
	pkghttp.JSON(w, http.StatusOK, h.service.List())
}
