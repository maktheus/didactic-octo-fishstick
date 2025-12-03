package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	bencrepo "github.com/example/back-end-tcc/services/benchmark/repository"
)

// Option customises benchmark service behaviour.
type Option func(*Service)

// WithLogger attaches a logger.
func WithLogger(l logger.Logger) Option {
	return func(s *Service) {
		s.log = l
	}
}

// WithMetrics attaches a metrics recorder.
func WithMetrics(rec metrics.Recorder) Option {
	return func(s *Service) {
		s.metrics = rec
	}
}

// Service manages benchmarks.
type Service struct {
	repo    *bencrepo.BenchmarkRepository
	log     logger.Logger
	metrics metrics.Recorder
}

// New creates service.
func New(repo *bencrepo.BenchmarkRepository, opts ...Option) *Service {
	svc := &Service{repo: repo, log: logger.New()}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// Create creates a benchmark definition.
// Create creates a benchmark definition.
func (s *Service) Create(b models.Benchmark) (models.Benchmark, error) {
	start := time.Now()
	if b.ID == "" {
		// Generate ID if missing
		b.ID = fmt.Sprintf("bench-%d", time.Now().UnixNano())
	}
	if b.Name == "" {
		s.observe("create", start, "error")
		return models.Benchmark{}, errors.New("missing name")
	}
	b.CreatedAt = time.Now()
	b.TasksCount = len(b.Tasks)

	s.repo.Save(b)
	if s.log != nil {
		s.log.Printf("benchmark: created benchmark %s", b.ID)
	}
	s.observe("create", start, "ok")
	return b, nil
}

// List returns benchmarks.
func (s *Service) List() []models.Benchmark {
	if s.metrics != nil {
		s.metrics.AddCounter("benchmark_list_total", map[string]string{"result": "ok"}, 1)
	}
	return s.repo.List()
}

func (s *Service) observe(operation string, start time.Time, result string) {
	if s.metrics == nil {
		return
	}
	labels := map[string]string{"operation": operation, "result": result}
	s.metrics.AddCounter("benchmark_operations_total", labels, 1)
	s.metrics.ObserveHistogram("benchmark_operation_duration_ms", labels, float64(time.Since(start).Milliseconds()))
}
