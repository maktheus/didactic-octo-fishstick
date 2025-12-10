package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// DockerSandbox implements Sandbox using Docker containers.
type DockerSandbox struct {
	cli         *client.Client
	containerID string
	image       string
	ctx         context.Context
}

// NewDockerSandbox creates a new DockerSandbox.
func NewDockerSandbox(image string) (*DockerSandbox, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.44"))
	if err != nil {
		return nil, err
	}
	return &DockerSandbox{
		cli:   cli,
		image: image,
		ctx:   context.Background(),
	}, nil
}

// Start creates and starts the container.
func (s *DockerSandbox) Start() error {
	// Ensure image exists (pull if needed)
	// For speed, we assume image exists or let CreateContainer fail if not found locally?
	// Better to pull if missing.
	_, _, err := s.cli.ImageInspectWithRaw(s.ctx, s.image)
	if client.IsErrNotFound(err) {
		reader, err := s.cli.ImagePull(s.ctx, s.image, types.ImagePullOptions{})
		if err != nil {
			return fmt.Errorf("failed to pull image: %w", err)
		}
		io.Copy(io.Discard, reader)
		reader.Close()
	}

	// Create container that stays alive
	resp, err := s.cli.ContainerCreate(s.ctx, &container.Config{
		Image: s.image,
		Cmd:   []string{"tail", "-f", "/dev/null"}, // Keep running
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	s.containerID = resp.ID

	if err := s.cli.ContainerStart(s.ctx, s.containerID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

// Stop stops and removes the container.
func (s *DockerSandbox) Stop() error {
	if s.containerID == "" {
		return nil
	}
	// Force remove
	return s.cli.ContainerRemove(s.ctx, s.containerID, types.ContainerRemoveOptions{Force: true})
}

// Exec executes a command inside the container.
func (s *DockerSandbox) Exec(cmd []string) (string, string, error) {
	if s.containerID == "" {
		return "", "", fmt.Errorf("sandbox not started")
	}

	execConfig := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	resp, err := s.cli.ContainerExecCreate(s.ctx, s.containerID, execConfig)
	if err != nil {
		return "", "", fmt.Errorf("failed to create exec: %w", err)
	}

	hijackedResp, err := s.cli.ContainerExecAttach(s.ctx, resp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", "", fmt.Errorf("failed to attach exec: %w", err)
	}
	defer hijackedResp.Close()

	var stdout, stderr bytes.Buffer
	// stdcopy.StdCopy demultiplexes the stream
	_, err = stdcopy.StdCopy(&stdout, &stderr, hijackedResp.Reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to copy output: %w", err)
	}

	// Wait for exec to finish to get exit code?
	// ContainerExecInspect can check exit code.
	for {
		inspect, err := s.cli.ContainerExecInspect(s.ctx, resp.ID)
		if err != nil {
			return stdout.String(), stderr.String(), err
		}
		if !inspect.Running {
			if inspect.ExitCode != 0 {
				return stdout.String(), stderr.String(), fmt.Errorf("command exited with code %d: %s", inspect.ExitCode, stderr.String())
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	return stdout.String(), stderr.String(), nil
}

// ID returns the container ID.
func (s *DockerSandbox) ID() string {
	return s.containerID
}
