package sandbox

import (
	"strings"
	"testing"
)

func TestDockerSandbox(t *testing.T) {
	// Skip if docker is not available?
	// For now, we assume it is since we are in verification.
	
	sb, err := NewDockerSandbox("python:3.9-slim")
	if err != nil {
		t.Fatalf("failed to create sandbox: %v", err)
	}

	if err := sb.Start(); err != nil {
		t.Fatalf("failed to start sandbox: %v", err)
	}
	defer sb.Stop()

	// Test 1: Simple Echo
	stdout, stderr, err := sb.Exec([]string{"echo", "hello"})
	if err != nil {
		t.Fatalf("exec failed: %v", err)
	}
	if strings.TrimSpace(stdout) != "hello" {
		t.Errorf("expected 'hello', got '%s'", stdout)
	}

	// Test 2: Write File and Read It
	_, _, err = sb.Exec([]string{"sh", "-c", "echo 'secret data' > /tmp/secret.txt"})
	if err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	stdout, stderr, err = sb.Exec([]string{"cat", "/tmp/secret.txt"})
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	if strings.TrimSpace(stdout) != "secret data" {
		t.Errorf("expected 'secret data', got '%s' (stderr: %s)", stdout, stderr)
	}
}
