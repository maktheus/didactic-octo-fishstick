package service

import (
	"context"
	"time"

	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/queue"
	scoringrepo "github.com/example/back-end-tcc/services/scoring/repository"
)

// ScoreCalculated queue message.
const ScoreCalculated = "score.calculated"

// Option customises the scoring service.
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

// Service aggregates metrics for submissions.
type Service struct {
	repo       *scoringrepo.ScoreRepository
	subscriber queue.Subscriber
	publisher  queue.Publisher
	log        logger.Logger
	metrics    metrics.Recorder
}

// New creates service.
func New(repo *scoringrepo.ScoreRepository, subscriber queue.Subscriber, publisher queue.Publisher, opts ...Option) *Service {
	s := &Service{repo: repo, subscriber: subscriber, publisher: publisher, log: logger.New()}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Start subscribes to runner results.
func (s *Service) Start() {
	s.subscriber.Subscribe(ScoreCalculated, s.handleSubmission)
	if s.log != nil {
		s.log.Println("scoring: subscribed to score.calculated")
	}
}

func (s *Service) handleSubmission(ctx context.Context, msg queue.Message) error {
	start := time.Now()
	submission, ok := msg.Data.(models.Submission)
	if !ok || submission.ScoreSummary == nil {
		s.observeScore(start, "ignored")
		return nil
	}
	summary := models.ScoreSummary{
		Score:      submission.ScoreSummary.Score,
		Metrics:    submission.ScoreSummary.Metrics,
		Calculated: submission.ScoreSummary.Calculated,
	}
	s.repo.Save(submission.ID, summary)
	if err := s.publisher.Publish(ctx, queue.Message{Type: "leaderboard.updated", Data: submission}); err != nil {
		if s.log != nil {
			s.log.Printf("scoring: failed to publish summary for submission %s: %v", submission.ID, err)
		}
		s.observeScore(start, "error")
		return err
	}
	if s.log != nil {
		s.log.Printf("scoring: stored summary for submission %s", submission.ID)
	}
	s.observeScore(start, "ok")
	return nil
}

// Summaries lists stored summaries.
func (s *Service) Summaries() []models.ScoreSummary {
	if s.metrics != nil {
		s.metrics.AddCounter("scoring_summaries_total", map[string]string{"result": "ok"}, float64(len(s.repo.List())))
	}
	return s.repo.List()
}

func (s *Service) observeScore(start time.Time, result string) {
	if s.metrics == nil {
		return
	}
	labels := map[string]string{"result": result}
	s.metrics.AddCounter("scoring_events_total", labels, 1)
	s.metrics.ObserveHistogram("scoring_duration_ms", labels, float64(time.Since(start).Milliseconds()))
}
