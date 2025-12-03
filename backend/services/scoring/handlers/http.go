package handlers

import (
	"net/http"

	pkghttp "github.com/example/back-end-tcc/pkg/http"
	"github.com/example/back-end-tcc/services/scoring/service"
)

// HTTP exposes scoring endpoints.
type HTTP struct {
	service *service.Service
}

// New creates handlers.
func New(service *service.Service) *HTTP {
	return &HTTP{service: service}
}

// List returns summaries.
func (h *HTTP) List(w http.ResponseWriter, r *http.Request) {
	pkghttp.JSON(w, http.StatusOK, h.service.Summaries())
}
