package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/queue"
	"github.com/example/back-end-tcc/pkg/sandbox"
	agentrepo "github.com/example/back-end-tcc/services/agent/repository"
	benchrepo "github.com/example/back-end-tcc/services/benchmark/repository"
	runnerrepo "github.com/example/back-end-tcc/services/runner/repository"
	"github.com/example/back-end-tcc/services/runner/patterns"
	"github.com/example/back-end-tcc/services/runner/tools"
)

// Option allows customizing service dependencies.
type Option func(*Service)

// WithLogger attaches a logger for instrumentation.
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

// Service consumes submissions and produces results.
type Service struct {
	repo          *runnerrepo.ResultRepository
	agentRepo     *agentrepo.AgentRepository
	benchmarkRepo *benchrepo.BenchmarkRepository
	subscriber    queue.Subscriber
	publisher     queue.Publisher
	log           logger.Logger
	metrics       metrics.Recorder
}

// New creates service.
func New(repo *runnerrepo.ResultRepository, agentRepo *agentrepo.AgentRepository, benchmarkRepo *benchrepo.BenchmarkRepository, subscriber queue.Subscriber, publisher queue.Publisher, opts ...Option) *Service {
	s := &Service{
		repo:          repo,
		agentRepo:     agentRepo,
		benchmarkRepo: benchmarkRepo,
		subscriber:    subscriber,
		publisher:     publisher,
		log:           logger.New(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Start registers queue consumers.
func (s *Service) Start() {
	s.subscriber.Subscribe("submission.created", s.handleSubmission)
	if s.log != nil {
		s.log.Println("runner: subscribed to submission.created")
	}
}

func (s *Service) handleSubmission(ctx context.Context, msg queue.Message) error {
	start := time.Now()
	submission, ok := msg.Data.(models.Submission)
	if !ok {
		s.observeRun(start, "ignored")
		return nil
	}
	if s.log != nil {
		s.log.Printf("runner: processing submission %s", submission.ID)
	}

	// Fetch Agent and Benchmark details
	agent, ok := s.agentRepo.Get(submission.AgentID)
	if !ok {
		s.log.Printf("runner: agent %s not found", submission.AgentID)
		return fmt.Errorf("agent not found")
	}
	benchmark, ok := s.benchmarkRepo.Get(submission.BenchmarkID)
	if !ok {
		s.log.Printf("runner: benchmark %s not found", submission.BenchmarkID)
		return fmt.Errorf("benchmark not found")
	}

	// Execute Tasks
	totalScore := 0.0
	tasksCount := len(benchmark.Tasks)
	if tasksCount == 0 {
		// Fallback if no tasks defined in benchmark
		tasksCount = 1
	}

	// For simplicity, we'll execute the first task or a generic prompt if no tasks
	prompt := "Hello, are you working?"
	if len(benchmark.Tasks) > 0 {
		prompt = benchmark.Tasks[0].Prompt
	}

	// Initialize Sandbox
	sb, err := sandbox.NewDockerSandbox("python:3.9-slim")
	if err != nil {
		s.log.Printf("runner: failed to create sandbox: %v", err)
		return err
	}
	if err := sb.Start(); err != nil {
		s.log.Printf("runner: failed to start sandbox: %v", err)
		return err
	}
	defer sb.Stop()

	// 1. Plan
	plan, err := patterns.GeneratePlan(&agent, prompt)
	if err != nil {
		s.log.Printf("runner: planning failed: %v", err)
		// Fallback to no plan? Or fail? Let's log and continue without plan for now, or fail.
		// For Level 4, let's just log.
	} else {
		s.log.Printf("runner: generated plan: %s", plan)
		// Inject plan into prompt
		prompt = fmt.Sprintf("Goal: %s\n\nPlan:\n%s\n\nExecute the plan using available tools.", prompt, plan)
		
		// Log Plan Trace
		s.repo.SaveTrace(models.TraceEvent{
			ID:           fmt.Sprintf("trace-%d", time.Now().UnixNano()),
			SubmissionID: submission.ID,
			Type:         "plan",
			Message:      plan,
			Timestamp:    time.Now(),
			Level:        "info",
		})
	}

	// 2. Execute & Reflect Loop
	maxRetries := 3
	var response string
	
	for i := 0; i < maxRetries; i++ {
		response, err = s.callOpenAI(&agent, prompt, sb)
		if err != nil {
			s.log.Printf("runner: execution failed: %v", err)
			submission.Status = "failed"
			break
		}

		// 3. Reflect
		approved, feedback, err := patterns.Reflect(&agent, benchmark.Tasks[0].Prompt, response)
		if err != nil {
			s.log.Printf("runner: reflection failed: %v", err)
			// If reflection fails, assume success or break?
			break
		}
		
		// Log Reflection Trace
		s.repo.SaveTrace(models.TraceEvent{
			ID:           fmt.Sprintf("trace-%d-reflect-%d", time.Now().UnixNano(), i),
			SubmissionID: submission.ID,
			Type:         "reflection",
			Message:      fmt.Sprintf("Approved: %v\nFeedback: %s", approved, feedback),
			Timestamp:    time.Now(),
			Level:        "info",
			Success:      approved,
		})

		if approved {
			s.log.Printf("runner: result approved")
			submission.Status = "completed"
			totalScore = 1.0
			break
		}

		s.log.Printf("runner: result rejected, retrying with feedback: %s", feedback)
		prompt = fmt.Sprintf("Previous attempt failed.\nFeedback: %s\n\nTry again.", feedback)
		// Continue loop
	}
	
	if submission.Status == "" {
		submission.Status = "failed" // If loop finishes without approval
	} else if submission.Status == "completed" {
		// Log final response
		s.log.Printf("runner: final response: %s", response)
	}

	now := time.Now()
	submission.CompletedAt = &now
	submission.ScoreSummary = &models.ScoreSummary{
		Score:      totalScore,
		Metrics:    map[string]float64{"accuracy": totalScore},
		Calculated: now,
	}

	// Important: Save to the shared repository
	s.repo.Save(submission)

	if err := s.publisher.Publish(ctx, queue.Message{Type: "score.calculated", Data: submission}); err != nil {
		if s.log != nil {
			s.log.Printf("runner: failed to publish score for submission %s: %v", submission.ID, err)
		}
		s.observeRun(start, "error")
		return err
	}
	if s.log != nil {
		s.log.Printf("runner: completed submission %s", submission.ID)
	}
	s.observeRun(start, "ok")
	return nil
}

func (s *Service) callOpenAI(agent *models.User, prompt string, sb sandbox.Sandbox) (string, error) {
	if agent.Model == "mock" {
		return s.mockLLM(prompt)
	}

	endpoint := agent.Endpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}

	model := agent.Model
	if model == "" {
		model = "gpt-4"
	}

	messages := []map[string]interface{}{}
	if agent.SystemPrompt != "" {
		messages = append(messages, map[string]interface{}{
			"role":    "system",
			"content": agent.SystemPrompt,
		})
	}
	messages = append(messages, map[string]interface{}{
		"role":    "user",
		"content": prompt,
	})

	availableTools := tools.GetTools()
	maxTurns := 10

	for i := 0; i < maxTurns; i++ {
		reqBody := map[string]interface{}{
			"model":    model,
			"messages": messages,
			"tools":    availableTools,
		}

		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return "", err
		}

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/json")
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("openai api error: %s - %s", resp.Status, string(body))
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return "", err
		}

		choices, ok := result["choices"].([]interface{})
		if !ok || len(choices) == 0 {
			return "", fmt.Errorf("no choices in response")
		}

		firstChoice := choices[0].(map[string]interface{})
		message := firstChoice["message"].(map[string]interface{})

		// Add assistant message to history
		messages = append(messages, message)

		// Check for tool calls
		if toolCalls, ok := message["tool_calls"].([]interface{}); ok && len(toolCalls) > 0 {
			s.log.Printf("runner: processing %d tool calls", len(toolCalls))
			for _, tc := range toolCalls {
				toolCall := tc.(map[string]interface{})
				function := toolCall["function"].(map[string]interface{})
				name := function["name"].(string)
				args := function["arguments"].(string)
				id := toolCall["id"].(string)

				s.log.Printf("runner: executing tool %s", name)
				output, err := tools.ExecuteTool(sb, name, args)
				if err != nil {
					output = fmt.Sprintf("Error executing tool: %v", err)
				}

				messages = append(messages, map[string]interface{}{
					"role":         "tool",
					"tool_call_id": id,
					"name":         name,
					"content":      output,
				})
			}
			continue // Loop again to send tool outputs to model
		}

		// No tool calls, return content
		if content, ok := message["content"].(string); ok {
			return content, nil
		}

		// If no content and no tool calls, something is wrong or it's just thinking
		return "", fmt.Errorf("no content or tool calls in response")
	}

	return "", fmt.Errorf("max turns reached")
}

func (s *Service) mockLLM(prompt string) (string, error) {
	s.log.Printf("runner: using mock LLM for prompt: %s", prompt)
	// Simple heuristic response
	if len(prompt) > 0 {
		return "Mock execution successful. I have completed the task.", nil
	}
	return "I am a mock agent.", nil
}

// Results returns processed submissions.
func (s *Service) Results() []models.Submission {
	if s.metrics != nil {
		s.metrics.AddCounter("runner_results_total", map[string]string{"result": "ok"}, float64(len(s.repo.List())))
	}
	return s.repo.List()
}

func (s *Service) observeRun(start time.Time, result string) {
	if s.metrics == nil {
		return
	}
	labels := map[string]string{"result": result}
	s.metrics.AddCounter("runner_runs_total", labels, 1)
	s.metrics.ObserveHistogram("runner_duration_ms", labels, float64(time.Since(start).Milliseconds()))
}
