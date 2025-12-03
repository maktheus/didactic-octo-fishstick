package patterns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/example/back-end-tcc/pkg/models"
)

// GeneratePlan calls the LLM to generate a plan for the given task.
func GeneratePlan(agent *models.User, taskPrompt string) (string, error) {
	if agent.Model == "mock" {
		return "1. Write python script\n2. Write test\n3. Run test", nil
	}

	endpoint := agent.Endpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}

	model := agent.Model
	if model == "" {
		model = "gpt-4"
	}

	systemPrompt := "You are an expert planner. Your goal is to break down a complex task into a clear, step-by-step execution plan. Do not execute the steps, just list them. Be concise."
	userPrompt := fmt.Sprintf("Task: %s\n\nCreate a numbered list of steps to complete this task.", taskPrompt)

	messages := []map[string]string{
		{"role": "system", "content": systemPrompt},
		{"role": "user", "content": userPrompt},
	}

	reqBody := map[string]interface{}{
		"model":    model,
		"messages": messages,
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

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("planner api error: %s - %s", resp.Status, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	firstChoice := choices[0].(map[string]interface{})
	message := firstChoice["message"].(map[string]interface{})
	content := message["content"].(string)

	return content, nil
}
