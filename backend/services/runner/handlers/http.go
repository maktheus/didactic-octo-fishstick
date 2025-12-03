package handlers

import (
	"net/http"

	pkghttp "github.com/example/back-end-tcc/pkg/http"
	"github.com/example/back-end-tcc/services/runner/service"
)

// HTTP exposes runner results.
type HTTP struct {
	service *service.Service
}

// New creates handlers.
func New(service *service.Service) *HTTP {
	return &HTTP{service: service}
}

// Results returns processed submissions.
func (h *HTTP) Results(w http.ResponseWriter, r *http.Request) {
	pkghttp.JSON(w, http.StatusOK, h.service.Results())
}
