package service

import (
	"context"
	"sort"
	"time"

	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/queue"
	lbrepo "github.com/example/back-end-tcc/services/leaderboard/repository"
)

// Option customises the leaderboard service.
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

// Service updates leaderboard projections.
type Service struct {
	repo    *lbrepo.Repository
	sub     queue.Subscriber
	log     logger.Logger
	metrics metrics.Recorder
}

// New creates service.
func New(repo *lbrepo.Repository, sub queue.Subscriber, opts ...Option) *Service {
	svc := &Service{repo: repo, sub: sub, log: logger.New()}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// Start listens for score events.
func (s *Service) Start() {
	s.sub.Subscribe("leaderboard.updated", s.handleScore)
	if s.log != nil {
		s.log.Println("leaderboard: subscribed to leaderboard.updated")
	}
}

func (s *Service) handleScore(ctx context.Context, msg queue.Message) error {
	start := time.Now()
	summary, ok := msg.Data.(models.ScoreSummary)
	if !ok {
		s.observe("ignored", start, 0)
		return nil
	}
	entry := models.LeaderboardEntry{
		SubmissionID: summary.Calculated.Format("20060102150405"),
		BenchmarkID:  "default",
		AgentID:      "agent",
		Score:        summary.Score,
	}
	entries := append(s.repo.List(), entry)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Score > entries[j].Score })
	for idx := range entries {
		entries[idx].Rank = idx + 1
		s.repo.Save(entries[idx])
	}
	if s.log != nil {
		s.log.Printf("leaderboard: updated with score %.2f", summary.Score)
	}
	s.observe("ok", start, summary.Score)
	return nil
}

// Entries returns leaderboard entries.
func (s *Service) Entries() []models.LeaderboardEntry {
	if s.metrics != nil {
		s.metrics.AddCounter("leaderboard_entries_total", map[string]string{"result": "ok"}, float64(len(s.repo.List())))
	}
	return s.repo.List()
}

func (s *Service) observe(result string, start time.Time, score float64) {
	if s.metrics == nil {
		return
	}
	labels := map[string]string{"result": result}
	s.metrics.AddCounter("leaderboard_events_total", labels, 1)
	s.metrics.ObserveHistogram("leaderboard_duration_ms", labels, float64(time.Since(start).Milliseconds()))
	if result == "ok" {
		s.metrics.ObserveHistogram("leaderboard_score_values", nil, score)
	}
}
