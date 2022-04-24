package container

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"time"
)

type ContainerCreateConfig struct {
	Image        string              // Name of the image as it was passed by the operator (e.g. could be symbolic)
	Entrypoint   ContainerCmdLine    // Entrypoint to run when starting the container
	Cmd          ContainerCmdLine    // Command to run when starting the container
	Env          []string            // List of environment variable to set in the container
	Volumes      map[string]struct{} // List of volumes (mounts) used for the container
	ExposedPorts nat.PortSet         `json:",omitempty"` // List of exposed ports
	Tty          bool                // Attach standard streams to a tty, including stdin if it is not closed.
	Links        []string            // List of links (in the name:alias form)
	NetworkMode  string              // Network mode to use for the container, e.g., --network=container:nginx
	Binds        []string            // List of volume bindings for this container
	PortBindings nat.PortMap         // Port mapping between the exposed port (container) and the host
	VolumesFrom  []string            // List of volumes to take from other container
}

type ContainerRemoveConfig = types.ContainerRemoveOptions

type ContainerListConfig = types.ContainerListOptions

type ContainerStartConfig = types.ContainerStartOptions

type ContainerStopConfig struct {
	timeout *time.Duration
}
