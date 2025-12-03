package sandbox

// Sandbox defines the interface for an isolated execution environment.
type Sandbox interface {
	// Start starts the sandbox.
	Start() error
	// Stop stops and cleans up the sandbox.
	Stop() error
	// Exec executes a command inside the sandbox and returns stdout, stderr, and error.
	Exec(cmd []string) (string, string, error)
	// ID returns the sandbox identifier (e.g., container ID).
	ID() string
}
