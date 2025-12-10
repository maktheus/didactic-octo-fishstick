package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/example/back-end-tcc/pkg/agent_protocol"
)

// AgentClient defines the interface for communicating with an agent.
type AgentClient interface {
	Step(url string, req agent_protocol.AgentRequest) (*agent_protocol.AgentResponse, error)
}

// HttpAgentClient implements AgentClient via HTTP.
type HttpAgentClient struct {
	client *http.Client
}

// NewHttpAgentClient creates a new HttpAgentClient.
func NewHttpAgentClient() *HttpAgentClient {
	return &HttpAgentClient{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Step sends a step request to the agent and returns the response.
func (c *HttpAgentClient) Step(url string, req agent_protocol.AgentRequest) (*agent_protocol.AgentResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.client.Post(url+"/agent/step", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to agent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("agent returned status %d", resp.StatusCode)
	}

	var agentResp agent_protocol.AgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&agentResp); err != nil {
		return nil, fmt.Errorf("failed to decode agent response: %w", err)
	}

	return &agentResp, nil
}
