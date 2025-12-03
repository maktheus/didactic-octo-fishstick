package http

import (
	"encoding/json"
	stdhttp "net/http"
)

// JSON writes the provided payload as JSON with the given status.
func JSON(w stdhttp.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

// Error writes an error response with a message.
func Error(w stdhttp.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}
