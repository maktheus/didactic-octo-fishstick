package models

import "time"

// User represents an authenticated subject or an Agent.
type User struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Email        string            `json:"email"`
	Role         string            `json:"role"`
	Provider     string            `json:"provider"`     // New: OpenAI, Anthropic, etc.
	Endpoint     string            `json:"endpoint"`     // New: API Endpoint
	Image        string            `json:"image"`        // New: Docker Image for Agent
	Model        string            `json:"model"`        // New: gpt-4, llama3, etc.
	SystemPrompt string            `json:"systemPrompt"` // New: Persona/Instructions
	AuthType     string            `json:"authType"`     // New: Bearer, API Key
	Status       string            `json:"status"`       // New: active, inactive
	CreatedAt    time.Time         `json:"createdAt"`    // New
	Headers      map[string]string `json:"headers"`      // New: Custom headers
	AuthToken    string            `json:"authToken"`    // New: Per-agent API Key
}

// Benchmark describes a benchmark suite definition.
type Benchmark struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Domain      string    `json:"domain"`     // New: Customer Support, Coding, etc.
	TasksCount  int       `json:"tasksCount"` // New
	Tasks       []Task    `json:"tasks"`      // New
	CreatedAt   time.Time `json:"createdAt"`
}

// Task describes a specific task within a benchmark.
type Task struct {
	ID             string   `json:"id"`
	Prompt         string   `json:"prompt"`
	ExpectedTool   string   `json:"expectedTool"`
	ExpectedOutput string   `json:"expected_output"` // New: For text matching
	Constraints    []string `json:"constraints"`
	MaxTurns       int      `json:"maxTurns"`
	// SWE-bench / Coding fields
	Repo       string   `json:"repo,omitempty"`
	Commit     string   `json:"commit,omitempty"`
	Patch      string   `json:"patch,omitempty"`
	TestFiles  []string `json:"testFiles,omitempty"`
	Difficulty string   `json:"difficulty,omitempty"`
}

// Submission is a benchmark submission by an agent (Run).
type Submission struct {
	ID            string        `json:"id"`
	AgentID       string        `json:"agentId"`
	AgentName     string        `json:"agentName"` // New: Denormalized for easier UI display
	BenchmarkID   string        `json:"benchmarkId"`
	BenchmarkName string        `json:"benchmarkName"` // New: Denormalized
	Payload       string        `json:"payload"`
	SubmittedAt   time.Time     `json:"submittedAt"`
	CompletedAt   *time.Time    `json:"completedAt"`
	Status        string        `json:"status"`
	Progress      int           `json:"progress"` // New: 0-100
	ScoreSummary  *ScoreSummary `json:"scoreSummary"`
}

// ScoreSummary captures scoring results.
type ScoreSummary struct {
	Score           float64            `json:"score"`
	SuccessRate     float64            `json:"successRate"`     // New
	ToolCorrectness float64            `json:"toolCorrectness"` // New
	Violations      int                `json:"violations"`      // New
	AvgTurns        float64            `json:"avgTurns"`        // New
	TotalCost       float64            `json:"totalCost"`       // New
	AvgLatency      float64            `json:"avgLatency"`      // New
	Metrics         map[string]float64 `json:"metrics"`
	Calculated      time.Time          `json:"calculated"`
}

// TraceEvent stores trace logs produced by benchmark runs.
type TraceEvent struct {
	ID           string            `json:"id"`
	SubmissionID string            `json:"submissionId"`
	TaskID       string            `json:"taskId"`   // New
	TaskName     string            `json:"taskName"` // New
	Type         string            `json:"type"`     // New: user, agent, tool
	Message      string            `json:"message"`  // Content
	ToolName     string            `json:"toolName"` // New
	Parameters   map[string]string `json:"parameters"`
	Result       map[string]string `json:"result"`
	Level        string            `json:"level"`
	Timestamp    time.Time         `json:"timestamp"`
	Success      bool              `json:"success"` // New
	Turns        int               `json:"turns"`   // New
	Cost         float64           `json:"cost"`    // New
	Latency      float64           `json:"latency"` // New
}

// LeaderboardEntry is a projection combining benchmark results.
type LeaderboardEntry struct {
	SubmissionID    string  `json:"submissionId"`
	BenchmarkID     string  `json:"benchmarkId"`
	AgentID         string  `json:"agentId"`
	AgentName       string  `json:"agentName"` // New
	Score           float64 `json:"score"`
	SuccessRate     float64 `json:"successRate"`     // New
	ToolCorrectness float64 `json:"toolCorrectness"` // New
	Violations      int     `json:"violations"`      // New
	AvgTurns        float64 `json:"avgTurns"`        // New
	TotalCost       float64 `json:"totalCost"`       // New
	AvgLatency      float64 `json:"avgLatency"`      // New
	Rank            int     `json:"rank"`
}
