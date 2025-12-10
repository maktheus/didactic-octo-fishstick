package patterns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/example/back-end-tcc/pkg/models"
)

// Reflect calls the LLM to critique the result of a task.
// Returns (approved, feedback).
func Reflect(agent *models.User, taskPrompt string, result string) (bool, string, error) {
	if agent.Model == "mock" {
		if strings.Contains(taskPrompt, "Failure") {
			return false, "REJECTED: Intentional failure for testing.", nil
		}
		if strings.Contains(taskPrompt, "infinite loop") {
			if strings.Contains(result, "Wrong Fix") {
				return false, "Still loops", nil
			}
			return true, "APPROVED: Loop fixed.", nil
		}
		return true, "APPROVED: Mock execution successful.", nil
	}

	endpoint := agent.Endpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}

	model := agent.Model
	if model == "" {
		model = "gpt-4"
	}

	systemPrompt := "You are a strict quality assurance engineer. Your goal is to verify if the result satisfies the original task. If it does, say 'APPROVED'. If not, explain what is missing or wrong."
	userPrompt := fmt.Sprintf("Original Task: %s\n\nResult:\n%s\n\nCritique:", taskPrompt, result)

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
		return false, "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, "", err
	}

	req.Header.Set("Content-Type", "application/json")
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, "", fmt.Errorf("reflector api error: %s - %s", resp.Status, string(body))
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, "", err
	}

	choices, ok := res["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return false, "", fmt.Errorf("no choices in response")
	}

	firstChoice := choices[0].(map[string]interface{})
	message := firstChoice["message"].(map[string]interface{})
	content := message["content"].(string)

	approved := strings.Contains(strings.ToUpper(content), "APPROVED")
	return approved, content, nil
}
