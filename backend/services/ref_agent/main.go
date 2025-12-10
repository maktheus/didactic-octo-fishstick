package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/example/back-end-tcc/pkg/agent_protocol"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/agent/step", handleAgentStep)

	log.Printf("Reference Agent listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleAgentStep(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req agent_protocol.AgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received step request for task %s, step %d", req.TaskID, req.StepNumber)

	// Simple logic: If step 1, list files. If step 2, finish.
	var resp agent_protocol.AgentResponse

	if req.StepNumber > 1 {
		resp = agent_protocol.AgentResponse{
			Thought: "I have seen enough. I am done.",
			Action:  "finish",
			ActionInput: map[string]interface{}{},
		}
	} else {
		resp = agent_protocol.AgentResponse{
			Thought: "I received a request. I will list the current directory to see where I am.",
			Action:  "run_command",
			ActionInput: map[string]string{
				"command": "ls -la",
			},
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
