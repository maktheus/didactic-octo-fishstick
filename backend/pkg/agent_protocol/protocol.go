package agent_protocol

// AgentRequest represents the payload sent to the agent's /agent/step endpoint.
type AgentRequest struct {
	TaskID      string                 `json:"task_id"`
	StepNumber  int                    `json:"step_number"`
	Input       string                 `json:"input"`
	Environment EnvironmentInfo        `json:"environment"`
	Tools       []ToolDefinition       `json:"tools"`
}

// EnvironmentInfo provides context about the execution environment.
type EnvironmentInfo struct {
	Cwd string `json:"cwd"`
	OS  string `json:"os"`
}

// ToolDefinition describes a tool available to the agent.
type ToolDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"` // JSON Schema as a map or struct
}

// AgentResponse represents the agent's reply to an execution step.
type AgentResponse struct {
	Thought     string      `json:"thought,omitempty"`
	Action      string      `json:"action"`
	ActionInput interface{} `json:"action_input"`
}
