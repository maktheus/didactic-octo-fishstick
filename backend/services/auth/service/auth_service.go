package service

import (
	"errors"
	"time"

	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	authrepo "github.com/example/back-end-tcc/services/auth/repository"
)

// Option customises the auth service.
type Option func(*AuthService)

// WithLogger attaches a logger.
func WithLogger(l logger.Logger) Option {
	return func(s *AuthService) {
		s.log = l
	}
}

// WithMetrics attaches a metrics recorder.
func WithMetrics(rec metrics.Recorder) Option {
	return func(s *AuthService) {
		s.metrics = rec
	}
}

// AuthService performs user authentication.
type AuthService struct {
	repo    *authrepo.UserRepository
	log     logger.Logger
	metrics metrics.Recorder
}

// NewAuthService creates a service.
func NewAuthService(repo *authrepo.UserRepository, opts ...Option) *AuthService {
	svc := &AuthService{repo: repo, log: logger.New()}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// Authenticate validates the supplied user identifier.
func (s *AuthService) Authenticate(id string) (models.User, error) {
	start := time.Now()
	if id == "" {
		s.observe(start, "error")
		return models.User{}, errors.New("missing user id")
	}
	user, ok := s.repo.FindByID(id)
	if !ok {
		s.observe(start, "error")
		return models.User{}, errors.New("user not found")
	}
	if s.log != nil {
		s.log.Printf("auth: authenticated user %s", id)
	}
	s.observe(start, "ok")
	return user, nil
}

func (s *AuthService) observe(start time.Time, result string) {
	if s.metrics == nil {
		return
	}
	labels := map[string]string{"result": result}
	s.metrics.AddCounter("auth_authenticate_total", labels, 1)
	s.metrics.ObserveHistogram("auth_authenticate_duration_ms", labels, float64(time.Since(start).Milliseconds()))
}
