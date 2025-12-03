package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/example/back-end-tcc/pkg/sandbox"
)

// ToolDefinition defines a tool for OpenAI.
type ToolDefinition struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function defines the function details.
type Function struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

// GetTools returns the list of available tools.
func GetTools() []ToolDefinition {
	return []ToolDefinition{
		{
			Type: "function",
			Function: Function{
				Name:        "read_file",
				Description: "Read the content of a file from the sandbox.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]string{
							"type":        "string",
							"description": "The absolute path to the file.",
						},
					},
					"required": []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name:        "write_file",
				Description: "Write content to a file in the sandbox. Overwrites if exists.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]string{
							"type":        "string",
							"description": "The absolute path to the file.",
						},
						"content": map[string]string{
							"type":        "string",
							"description": "The content to write.",
						},
					},
					"required": []string{"path", "content"},
				},
			},
		},
		{
			Type: "function",
			Function: Function{
				Name:        "run_command",
				Description: "Execute a shell command in the sandbox.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"command": map[string]string{
							"type":        "string",
							"description": "The command to execute.",
						},
					},
					"required": []string{"command"},
				},
			},
		},
	}
}

// ExecuteTool executes a tool in the sandbox.
func ExecuteTool(sb sandbox.Sandbox, name string, args string) (string, error) {
	var arguments map[string]interface{}
	if err := json.Unmarshal([]byte(args), &arguments); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	switch name {
	case "read_file":
		path, ok := arguments["path"].(string)
		if !ok {
			return "", fmt.Errorf("missing path")
		}
		// Use cat to read file
		stdout, stderr, err := sb.Exec([]string{"cat", path})
		if err != nil {
			return fmt.Sprintf("Error: %v\nStderr: %s", err, stderr), nil
		}
		if stderr != "" {
			return fmt.Sprintf("Stderr: %s", stderr), nil
		}
		return stdout, nil

	case "write_file":
		path, ok := arguments["path"].(string)
		if !ok {
			return "", fmt.Errorf("missing path")
		}
		content, ok := arguments["content"].(string)
		if !ok {
			return "", fmt.Errorf("missing content")
		}
		// Use sh -c to write file. Need to be careful with quoting.
		// A safer way is to write to a temp file via Exec input if supported, but Docker Exec doesn't easily support stdin.
		// For MVP, we'll try simple echo or printf. printf is safer.
		// escape single quotes
		escapedContent := strings.ReplaceAll(content, "'", "'\\''")
		cmd := fmt.Sprintf("printf '%%s' '%s' > '%s'", escapedContent, path)
		stdout, stderr, err := sb.Exec([]string{"sh", "-c", cmd})
		if err != nil {
			return fmt.Sprintf("Error: %v\nStderr: %s", err, stderr), nil
		}
		if stderr != "" {
			return fmt.Sprintf("Stderr: %s", stderr), nil
		}
		return fmt.Sprintf("Successfully wrote to %s\nOutput: %s", path, stdout), nil

	case "run_command":
		command, ok := arguments["command"].(string)
		if !ok {
			return "", fmt.Errorf("missing command")
		}
		stdout, stderr, err := sb.Exec([]string{"sh", "-c", command})
		if err != nil {
			return fmt.Sprintf("Error: %v\nStderr: %s", err, stderr), nil
		}
		return fmt.Sprintf("Stdout:\n%s\nStderr:\n%s", stdout, stderr), nil

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}
