package service

import (
	"fmt"
	"time"

	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	agentrepo "github.com/example/back-end-tcc/services/agent/repository"
)

// Option customises the agent service.
type Option func(*AgentService)

// WithLogger attaches a logger.
func WithLogger(l logger.Logger) Option {
	return func(s *AgentService) {
		s.log = l
	}
}

// WithMetrics attaches a recorder.
func WithMetrics(rec metrics.Recorder) Option {
	return func(s *AgentService) {
		s.metrics = rec
	}
}

// AgentService manages benchmark agents.
type AgentService struct {
	repo    *agentrepo.AgentRepository
	log     logger.Logger
	metrics metrics.Recorder
}

// NewAgentService creates service.
func NewAgentService(repo *agentrepo.AgentRepository, opts ...Option) *AgentService {
	svc := &AgentService{repo: repo, log: logger.New()}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// Register registers an agent.
func (s *AgentService) Register(agent *models.User) error {
	start := time.Now()
	if s.log != nil {
		s.log.Printf("agent: registering agent, initial ID: %s", agent.ID)
	}
	if agent.ID == "" {
		agent.ID = fmt.Sprintf("agent-%d", time.Now().UnixNano())
		if s.log != nil {
			s.log.Printf("agent: generated ID: %s", agent.ID)
		}
	}
	if agent.CreatedAt.IsZero() {
		agent.CreatedAt = time.Now()
	}
	if agent.Status == "" {
		agent.Status = "active"
	}
	s.repo.Save(*agent)
	if s.log != nil {
		s.log.Printf("agent: registered agent %s", agent.ID)
	}
	s.observe("register", start, "ok")
	return nil
}

// List agents.
func (s *AgentService) List() []models.User {
	if s.metrics != nil {
		s.metrics.AddCounter("agent_list_total", map[string]string{"result": "ok"}, 1)
	}
	return s.repo.List()
}

func (s *AgentService) observe(operation string, start time.Time, result string) {
	if s.metrics == nil {
		return
	}
	labels := map[string]string{"operation": operation, "result": result}
	s.metrics.AddCounter("agent_operations_total", labels, 1)
	s.metrics.ObserveHistogram("agent_operation_duration_ms", labels, float64(time.Since(start).Milliseconds()))
}
