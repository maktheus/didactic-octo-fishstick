package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const (
	baseURL = "http://localhost:8080"
)

// Minimal structs for JSON encoding/decoding
type Agent struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Endpoint string `json:"endpoint"`
	Image    string `json:"image"`
	Model    string `json:"model"`
}

type Benchmark struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Tasks []struct {
		Prompt string `json:"prompt"`
	} `json:"tasks"`
}

type Submission struct {
	BenchmarkID string `json:"benchmark_id"`
	AgentID     string `json:"agent_id"`
	Payload     string `json:"payload"`
}

type RunStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func TestRealAgentProtocolFlow(t *testing.T) {
	// 1. Health Check
	if err := waitForBackend(); err != nil {
		t.Fatalf("Backend not ready: %v", err)
	}

	// 2. Create Agent
	agent := Agent{
		Name:     "E2E Ref Agent",
		Provider: "agent-protocol",
		Endpoint: "http://ref-agent:8080", // Internal Docker DNS
		Image:    "",                        // Empty to bypass dynamic creation, using pre-existing service
		Model:    "mock",                    // "mock" triggers local reflection logic
	}
	agentID, err := postResource("/agents", agent)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}
	fmt.Printf("Created Agent: %s\n", agentID)

	// 3. Create Benchmark
	bench := Benchmark{
		Name: "E2E Benchmark",
		Tasks: []struct {
			Prompt string `json:"prompt"`
		}{
			{Prompt: "Run ls -la"},
		},
	}
	benchID, err := postResource("/benchmarks", bench)
	if err != nil {
		t.Fatalf("Failed to create benchmark: %v", err)
	}
	fmt.Printf("Created Benchmark: %s\n", benchID)

	// 4. Submit Run
	sub := Submission{
		BenchmarkID: benchID,
		AgentID:     agentID,
		Payload:     "{}",
	}
	runID, err := postResource("/submissions", sub)
	if err != nil {
		t.Fatalf("Failed to create submission: %v", err)
	}
	fmt.Printf("Started Run: %s\n", runID)

	// 5. Poll for Completion
	timeout := time.After(60 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Test timed out waiting for run completion")
		case <-ticker.C:
			status, err := getRunStatus(runID)
			if err != nil {
				t.Logf("Error getting status: %v", err)
				continue
			}
			fmt.Printf("Run Status: %s\n", status)
			if status == "completed" {
				return // Success!
			}
			if status == "failed" {
				t.Fatal("Run failed")
			}
		}
	}
}

func waitForBackend() error {
	for i := 0; i < 10; i++ {
		resp, err := http.Get(baseURL + "/health")
		if err == nil && (resp.StatusCode == 200 || resp.StatusCode == 404) { // 404 on health might happen if not registered, but connection worked
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("backend unreachable")
}

func postResource(path string, data interface{}) (string, error) {
	b, _ := json.Marshal(data)
	resp, err := http.Post(baseURL+path, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad status: %d body: %s", resp.StatusCode, string(body))
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if id, ok := res["id"].(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("no id returned")
}

func getRunStatus(id string) (string, error) {
	resp, err := http.Get(baseURL + "/submissions")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var runs []RunStatus
	if err := json.NewDecoder(resp.Body).Decode(&runs); err != nil {
		return "", err
	}

	for _, r := range runs {
		if r.ID == id {
			return r.Status, nil
		}
	}
	return "", fmt.Errorf("run not found")
}
