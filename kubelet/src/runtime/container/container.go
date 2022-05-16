package container

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"minik8s/kubelet/src/runtime/docker"
	"minik8s/util/logger"
	"time"
)

type State string

const (
	// StateCreated indicates a container that has been created (e.g. with docker create) but not started.
	StateCreated State = "created"
	// StateRunning indicates a currently running container.
	StateRunning State = "running"
	// StateExited indicates a container that ran and completed ("stopped" in other contexts, although a created container is technically also "stopped").
	StateExited State = "exited"
	// StateUnknown encompasses all the states that we currently don't care about (like restarting, paused, dead).
	StateUnknown State = "unknown"
)

// Status represents the status of a container.
type Status struct {
	// ID of the container.
	ID ID
	// Name of the container.
	Name string
	// Status of the container.
	State State
	// Creation time of the container.
	CreatedAt time.Time
	// Start time of the container.
	StartedAt time.Time
	// Finish time of the container.
	FinishedAt time.Time
	// Exit code of the container.
	ExitCode int
	// ID of the image.
	ImageID string
	// Number of times that the container has been restarted.
	RestartCount int
	// A string stands for the error
	Error string
	// The status of resource usage
	ResourcesUsage ResourcesUsage
}

type ID = string

type CmdLine = []string

type Container struct {
	ID      ID
	Name    string
	Image   string
	ImageID string
	State   State
}

type ResourcesUsage struct {
	CpuPercent float64 `json:"cpu_percent"`
	MemPercent float64 `json:"mem_percent"`
}

type Manager interface {
	ListContainers(config *ListConfig) ([]*Container, error)
	CreateContainer(name string, config *ContainerCreateConfig) (string, error)
	RemoveContainer(ID ID, config *RemoveConfig) error
	StartContainer(ID ID, config *StartConfig) error
	RenameContainer(ID ID, newName string) error
	StopContainer(ID ID, config *StopConfig) error
	GetContainerStats(ID ID) (ResourcesUsage, error)
	InspectContainer(ID ID) (InspectInfo, error)
}

func NewContainerManager() Manager {
	return &containerManager{}
}

type containerManager struct {
}

func (cm *containerManager) GetContainerStats(ID ID) (ResourcesUsage, error) {
	ru := ResourcesUsage{}

	resp, err := docker.Client.ContainerStats(docker.Ctx, ID, false)
	if err != nil {
		return ru, err
	}

	body := resp.Body
	osType := resp.OSType
	defer body.Close()
	decoder := json.NewDecoder(body)
	statsJson := &types.StatsJSON{}
	err = decoder.Decode(statsJson)
	if err != nil {
		return ru, err
	}

	if osType != "windows" {
		ru.MemPercent = calculateMemPercentUnix(statsJson)
		ru.CpuPercent = calculateCPUPercentUnix(statsJson)
	} else {
		ru.CpuPercent = calculateCPUPercentWindows(statsJson)
	}
	return ru, nil
}

func (cm *containerManager) RenameContainer(ID ID, newName string) error {
	return docker.Client.ContainerRename(docker.Ctx, ID, newName)
}

func (cm *containerManager) ListContainers(config *ListConfig) ([]*Container, error) {
	filter := filters.NewArgs()
	for name, value := range config.LabelSelector {
		if len(name) == 0 {
			continue
		}
		if len(value) > 0 {
			filter.Add("label", name+"="+value)
		} else { // filter just by name
			filter.Add("label", name)
		}
	}

	containers, err := docker.Client.ContainerList(docker.Ctx, types.ContainerListOptions{
		Quiet:   config.Quiet,
		Size:    config.Size,
		All:     config.All,
		Latest:  config.Latest,
		Since:   config.Since,
		Before:  config.Before,
		Limit:   config.Limit,
		Filters: filter,
	})
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
			State:   State(c.State),
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
		Labels:       config.Labels,
	}, &container.HostConfig{
		Binds:        config.Binds,
		PortBindings: config.PortBindings,
		NetworkMode:  config.NetworkMode,
		PidMode:      config.PidMode,
		IpcMode:      config.IpcMode,
		VolumesFrom:  config.VolumesFrom,
		Links:        config.Links,
	}, nil, nil, name)
	for _, warning := range res.Warnings {
		logger.Warn(warning)
	}

	return res.ID, err
}

func (cm *containerManager) RemoveContainerByName(name string, config *RemoveConfig) error {
	return docker.Client.ContainerRemove(docker.Ctx, name, *config)
}

func (cm *containerManager) RemoveContainer(ID ID, config *RemoveConfig) error {
	return docker.Client.ContainerRemove(docker.Ctx, ID, *config)
}

func (cm *containerManager) StartContainer(ID ID, config *StartConfig) error {
	return docker.Client.ContainerStart(docker.Ctx, ID, *config)
}

func (cm *containerManager) StopContainer(ID ID, config *StopConfig) error {
	return docker.Client.ContainerStop(docker.Ctx, ID, &config.timeout)
}

func (cm *containerManager) InspectContainer(ID ID) (InspectInfo, error) {
	return docker.Client.ContainerInspect(docker.Ctx, ID)
}
