package container

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"testDocker/kubelet/src/runtime/docker"
)

type ContainerState string

const (
	// ContainerStateCreated indicates a container that has been created (e.g. with docker create) but not started.
	ContainerStateCreated ContainerState = "created"
	// ContainerStateRunning indicates a currently running container.
	ContainerStateRunning ContainerState = "running"
	// ContainerStateExited indicates a container that ran and completed ("stopped" in other contexts, although a created container is technically also "stopped").
	ContainerStateExited ContainerState = "exited"
	// ContainerStateUnknown encompasses all the states that we currently don't care about (like restarting, paused, dead).
	ContainerStateUnknown ContainerState = "unknown"
)

type ContainerID = string

type ContainerCmdLine = []string

type Container struct {
	ID      ContainerID
	Name    string
	Image   string
	ImageID string
	State   ContainerState
}

type ContainerManager interface {
	ListContainer(config *ContainerListConfig) ([]*Container, error)
	CreateContainer(name string, config *ContainerCreateConfig) (string, error)
	RemoveContainer(ID ContainerID, config *ContainerRemoveConfig) error
	StartContainer(ID ContainerID, config *ContainerStartConfig) error
	StopContainer(ID ContainerID, config *ContainerStopConfig) error
	InspectContainer(ID ContainerID) (ContainerInspectInfo, error)
}

func NewContainerManager() ContainerManager {
	return &containerManager{}
}

type containerManager struct {
}

func (cm *containerManager) ListContainer(config *ContainerListConfig) ([]*Container, error) {
	containers, err := docker.Client.ContainerList(docker.Ctx, *config)
	if err != nil {
		return nil, err
	}

	ret := make([]*Container, len(containers))
	for i, c := range containers {
		ret[i] = &Container{
			ID:      c.ID,
			Name:    c.Names[0],
			Image:   c.Image,
			ImageID: c.ImageID,
			State:   ContainerState(c.State),
		}
	}
	return ret, nil
}

func (cm *containerManager) CreateContainer(name string, config *ContainerCreateConfig) (string, error) {
	res, err := docker.Client.ContainerCreate(docker.Ctx, &container.Config{
		ExposedPorts: config.ExposedPorts,
		Tty:          config.Tty,
		Env:          config.Env,
		Cmd:          config.Cmd,
		Image:        config.Image,
		Volumes:      config.Volumes,
		Entrypoint:   config.Entrypoint,
	}, &container.HostConfig{
		Binds:        config.Binds,
		PortBindings: config.PortBindings,
		NetworkMode:  container.NetworkMode(config.NetworkMode),
		VolumesFrom:  config.VolumesFrom,
		Links:        config.Links,
	}, nil, nil, name)

	for _, warning := range res.Warnings {
		fmt.Println(warning)
	}

	return res.ID, err
}

func (cm *containerManager) RemoveContainer(ID ContainerID, config *ContainerRemoveConfig) error {
	return docker.Client.ContainerRemove(docker.Ctx, ID, *config)
}

func (cm *containerManager) StartContainer(ID ContainerID, config *ContainerStartConfig) error {
	return docker.Client.ContainerStart(docker.Ctx, ID, *config)
}

func (cm *containerManager) StopContainer(ID ContainerID, config *ContainerStopConfig) error {
	return docker.Client.ContainerStop(docker.Ctx, ID, &config.timeout)
}

func (cm *containerManager) InspectContainer(ID ContainerID) (ContainerInspectInfo, error) {
	return docker.Client.ContainerInspect(docker.Ctx, ID)
}
