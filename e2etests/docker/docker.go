package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func RunContainer(containerName, networkName, imageName, configPath, certsPath string, ports map[string]string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Setup port bindings and exposed ports
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}
	for containerPort, hostPort := range ports {
		port, err := nat.NewPort("tcp", containerPort)
		if err != nil {
			return fmt.Errorf("invalid port %s: %v", containerPort, err)
		}
		portBindings[port] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: hostPort}}
		exposedPorts[port] = struct{}{}
	}

	// Include ExposedPorts in the container configuration
	contConfig := &container.Config{
		Image:        imageName,
		ExposedPorts: exposedPorts, // Explicitly declare exposed ports
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		NetworkMode:  container.NetworkMode(networkName),
		Binds: []string{
			fmt.Sprintf("%s:/etc/sepp/config.yaml", configPath),
			fmt.Sprintf("%s:/e2etests/certs/", certsPath),
		},
	}

	resp, err := cli.ContainerCreate(ctx, contConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	log.Printf("docker - container %s started successfully with exposed ports", containerName)
	return nil
}

func StopAndRemoveContainer(containerID string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	ctx := context.Background()

	stopOptions := container.StopOptions{}
	if err := cli.ContainerStop(ctx, containerID, stopOptions); err != nil {
		log.Printf("Warning: Failed to stop container %s: %v", containerID, err)
	}

	removeOptions := container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}
	if err := cli.ContainerRemove(ctx, containerID, removeOptions); err != nil {
		return fmt.Errorf("failed to remove container %s: %v", containerID, err)
	}
	log.Printf("docker - container %s removed successfully", containerID)
	return nil
}

func CreateNetwork(networkName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	_, err = cli.NetworkCreate(context.Background(), networkName, types.NetworkCreate{})
	log.Printf("docker - network %s created successfully", networkName)
	return err
}

func RemoveNetwork(networkName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	ctx := context.Background()

	if err := cli.NetworkRemove(ctx, networkName); err != nil {
		return fmt.Errorf("failed to remove network %s: %v", networkName, err)
	}

	log.Printf("docker - network %s removed successfully", networkName)
	return nil
}
