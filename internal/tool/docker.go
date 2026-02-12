package tool

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	dockerclient "github.com/docker/docker/client"
)

// getDockerClient creates a new Docker client
func getDockerClient() (*dockerclient.Client, error) {
	dockerAPIVersion := os.Getenv("DOCKER_API_VERSION")
	if dockerAPIVersion == "" {
		return dockerclient.NewClientWithOpts(dockerclient.FromEnv)
	}
	return dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithVersion(dockerAPIVersion))
}

// StartContainer starts a Docker container by name or ID
func StartContainer(containerName string) error {
	cli, err := getDockerClient()
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	defer cli.Close()

	ctx := context.Background()

	// Check if container exists and get its status
	containerJSON, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		return fmt.Errorf("failed to inspect container '%s': %w", containerName, err)
	}

	// Check if already running
	if containerJSON.State.Running {
		return fmt.Errorf("container '%s' is already running", containerName)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, containerName, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container '%s': %w", containerName, err)
	}

	return nil
}
