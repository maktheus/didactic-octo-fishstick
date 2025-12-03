package service

import (
	"context"
	"fmt"
	"time"

	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/queue"
	tracerepo "github.com/example/back-end-tcc/services/trace/repository"
)

// TraceCreated queue message type.
const TraceCreated = "trace.created"

// Option customises trace service.
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

// Service persists trace events and publishes notifications.
type Service struct {
	repo    *tracerepo.Repository
	bus     queue.Publisher
	log     logger.Logger
	metrics metrics.Recorder
}

// New creates service.
func New(repo *tracerepo.Repository, bus queue.Publisher, opts ...Option) *Service {
	svc := &Service{repo: repo, bus: bus, log: logger.New()}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// Record stores a trace event.
func (s *Service) Record(ctx context.Context, submissionID, message, level string) models.TraceEvent {
	start := time.Now()
	event := models.TraceEvent{
		ID:           generateTraceID(),
		SubmissionID: submissionID,
		Message:      message,
		Level:        level,
		Timestamp:    time.Now(),
	}
	s.repo.Save(event)
	if err := s.bus.Publish(ctx, queue.Message{Type: TraceCreated, Data: event}); err != nil {
		if s.log != nil {
			s.log.Printf("trace: failed to publish event %s: %v", event.ID, err)
		}
		s.observe("error", start)
		return event
	}
	if s.log != nil {
		s.log.Printf("trace: recorded event %s for submission %s", event.ID, submissionID)
	}
	s.observe("ok", start)
	return event
}

// Events returns stored traces.
func (s *Service) Events() []models.TraceEvent {
	if s.metrics != nil {
		s.metrics.AddCounter("trace_events_list_total", map[string]string{"result": "ok"}, 1)
	}
	return s.repo.List()
}

func (s *Service) observe(result string, start time.Time) {
	if s.metrics == nil {
		return
	}
	labels := map[string]string{"result": result}
	s.metrics.AddCounter("trace_record_total", labels, 1)
	s.metrics.ObserveHistogram("trace_record_duration_ms", labels, float64(time.Since(start).Milliseconds()))
}

func generateTraceID() string {
	return fmt.Sprintf("trace-%d", time.Now().UnixNano())
}
