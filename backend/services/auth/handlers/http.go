package handlers

import (
	"net/http"

	pkghttp "github.com/example/back-end-tcc/pkg/http"
	"github.com/example/back-end-tcc/services/auth/service"
)

// HTTP exposes auth HTTP handlers.
type HTTP struct {
	service *service.AuthService
}

// NewHTTP creates handlers.
func NewHTTP(service *service.AuthService) *HTTP {
	return &HTTP{service: service}
}

// Authenticate handles authentication requests.
func (h *HTTP) Authenticate(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("user_id")
	user, err := h.service.Authenticate(id)
	if err != nil {
		pkghttp.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	pkghttp.JSON(w, http.StatusOK, user)
}
