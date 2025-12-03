package handlers

import (
	"encoding/json"
	"net/http"

	pkghttp "github.com/example/back-end-tcc/pkg/http"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/services/benchmark/service"
)

// HTTP handles benchmark endpoints.
type HTTP struct {
	service *service.Service
}

// New creates handlers.
func New(service *service.Service) *HTTP {
	return &HTTP{service: service}
}

// Create handles creation requests.
func (h *HTTP) Create(w http.ResponseWriter, r *http.Request) {
	var benchmark models.Benchmark
	if err := json.NewDecoder(r.Body).Decode(&benchmark); err != nil {
		pkghttp.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}
	created, err := h.service.Create(benchmark)
	if err != nil {
		pkghttp.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	pkghttp.JSON(w, http.StatusCreated, created)
}

// List returns benchmarks.
func (h *HTTP) List(w http.ResponseWriter, r *http.Request) {
	pkghttp.JSON(w, http.StatusOK, h.service.List())
}
