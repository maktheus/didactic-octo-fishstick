package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/queue"
	orchrepo "github.com/example/back-end-tcc/services/orchestrator/repository"
)

// SubmissionCreated message type.
const SubmissionCreated = "submission.created"

// Option customises the service dependencies.
type Option func(*Service)

// WithLogger attaches a logger to the service.
func WithLogger(l logger.Logger) Option {
	return func(s *Service) {
		s.log = l
	}
}

// WithMetrics attaches a recorder for service telemetry.
func WithMetrics(rec metrics.Recorder) Option {
	return func(s *Service) {
		s.metrics = rec
	}
}

// Service coordinates benchmark submissions.
type Service struct {
	repo    *orchrepo.SubmissionRepository
	bus     queue.Publisher
	log     logger.Logger
	metrics metrics.Recorder
}

// New creates service.
func New(repo *orchrepo.SubmissionRepository, bus queue.Publisher, opts ...Option) *Service {
	s := &Service{repo: repo, bus: bus, log: logger.New()}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Submit registers a submission and notifies workers.
func (s *Service) Submit(ctx context.Context, benchmarkID, agentID, payload string) (models.Submission, error) {
	start := time.Now()
	if benchmarkID == "" || agentID == "" {
		s.observeSubmit(start, "error")
		return models.Submission{}, errors.New("missing identifiers")
	}
	submission := models.Submission{
		ID:          generateSubmissionID(),
		BenchmarkID: benchmarkID,
		AgentID:     agentID,
		Payload:     payload,
		SubmittedAt: time.Now(),
		Status:      "queued",
	}
	s.repo.Save(submission)
	if s.log != nil {
		s.log.Printf("orchestrator: submission %s queued for benchmark=%s agent=%s", submission.ID, benchmarkID, agentID)
	}
	if err := s.bus.Publish(ctx, queue.Message{Type: SubmissionCreated, Data: submission}); err != nil {
		if s.log != nil {
			s.log.Printf("orchestrator: failed to publish submission %s: %v", submission.ID, err)
		}
		s.observeSubmit(start, "error")
		return models.Submission{}, err
	}
	s.observeSubmit(start, "ok")
	return submission, nil
}

// List returns submissions.
func (s *Service) List() []models.Submission {
	if s.metrics != nil {
		s.metrics.AddCounter("orchestrator_list_total", map[string]string{"result": "ok"}, 1)
	}
	return s.repo.List()
}

func (s *Service) observeSubmit(start time.Time, result string) {
	if s.metrics == nil {
		return
	}
	labels := map[string]string{"result": result}
	s.metrics.AddCounter("orchestrator_submit_total", labels, 1)
	s.metrics.ObserveHistogram("orchestrator_submit_duration_ms", labels, float64(time.Since(start).Milliseconds()))
}

func generateSubmissionID() string {
	return fmt.Sprintf("sub-%d", time.Now().UnixNano())
}
