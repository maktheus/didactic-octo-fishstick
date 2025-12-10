package agent_container

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Manager handles the lifecycle of agent containers.
type Manager struct {
	cli *client.Client
}

// NewManager creates a new Manager.
func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.44"))
	if err != nil {
		return nil, err
	}
	return &Manager{cli: cli}, nil
}

// StartAgentContainer starts an agent container and returns its IP/Host and cleanup function.
// For simplicity, we expose port 8080 to a random host port and communicate via host mapping,
// OR we can use a shared network. For now, Host Mapping is easier for the Runner (running on host or in container?).
// If Runner is in a container, we should use a shared network.
// If Runner is on host (go run), we use port mapping.
// Let's assume Port Mapping for dev, but Shared Network is better for prod.
// We will try to inspect the container to get its IP if on same network, or map port.
func (m *Manager) StartAgentContainer(ctx context.Context, image string, envVars map[string]string) (string, func(), error) {
	// Pull image
	_, _, err := m.cli.ImageInspectWithRaw(ctx, image)
	if client.IsErrNotFound(err) {
		reader, err := m.cli.ImagePull(ctx, image, types.ImagePullOptions{})
		if err != nil {
			return "", nil, fmt.Errorf("failed to pull image: %w", err)
		}
		defer reader.Close()
		// Drain reader
		buffer := make([]byte, 1024)
		for {
			if _, err := reader.Read(buffer); err != nil {
				break
			}
		}
	}

	// Configure container
	config := &container.Config{
		Image: image,
		Env:   mapToEnvList(envVars),
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "0", // Random available port
				},
			},
		},
		AutoRemove: true,
	}

	resp, err := m.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create container: %w", err)
	}

	if err := m.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", nil, fmt.Errorf("failed to start container: %w", err)
	}

	cleanup := func() {
		timeout := 0 // Force kill immediately if needed
		// Use a new context for cleanup in case original is cancelled
		m.cli.ContainerStop(context.Background(), resp.ID, container.StopOptions{Timeout: &timeout})
	}

	// Inspect to find the mapped port
	inspect, err := m.cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	bindings := inspect.NetworkSettings.Ports["8080/tcp"]
	if len(bindings) == 0 {
		cleanup()
		return "", nil, fmt.Errorf("no port bindings found")
	}

	hostPort := bindings[0].HostPort
	endpoint := fmt.Sprintf("http://localhost:%s", hostPort)

	// Wait for health check?
	// The runner loop will fail if it can't connect, so maybe wait a bit here or let retry logic handle it.
	time.Sleep(2 * time.Second) 

	return endpoint, cleanup, nil
}

func mapToEnvList(m map[string]string) []string {
	var list []string
	for k, v := range m {
		list = append(list, fmt.Sprintf("%s=%s", k, v))
	}
	return list
}
