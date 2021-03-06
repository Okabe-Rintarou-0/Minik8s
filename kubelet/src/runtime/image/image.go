package image

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"minik8s/kubelet/src/runtime/docker"
	"minik8s/util/logger"
)

var log = logger.Log("Image puller")

type Summary = types.ImageSummary

type Manager interface {
	PullImage(name string, config *PullConfig) error
	ExistsImage(name string) (bool, error)
	ListImages(config *ListConfig) ([]Summary, error)
	ListImagesByName(name string, config *ListConfig) ([]Summary, error)
	RemoveImage(ID string, config *RemoveConfig) (RemoveResponse, error)
}

type imageManager struct {
}

func (is *imageManager) ExistsImage(name string) (bool, error) {
	images, err := is.ListImagesByName(name, &ListConfig{All: true})
	return len(images) > 0, err
}

func (is *imageManager) ListImages(config *ListConfig) ([]Summary, error) {
	return docker.Client.ImageList(docker.Ctx, types.ImageListOptions{
		All: config.All,
	})
}

func (is *imageManager) ListImagesByName(name string, config *ListConfig) ([]Summary, error) {
	filter := filters.NewArgs()
	filter.Add("dangling", "false")
	filter.Add("reference", name)
	return docker.Client.ImageList(docker.Ctx, types.ImageListOptions{
		All:     config.All,
		Filters: filter,
	})
}

func (is *imageManager) PullImage(name string, config *PullConfig) error {
	events, err := docker.Client.ImagePull(docker.Ctx, name, types.ImagePullOptions{
		All: config.All,
	})
	if err != nil {
		return err
	}
	if config.Verbose {
		parseAndPrintPullEvents(events, name)
	}
	return nil
}

func (is *imageManager) RemoveImage(ID string, config *RemoveConfig) (RemoveResponse, error) {
	resp, err := docker.Client.ImageRemove(docker.Ctx, ID, *config)
	return RemoveResponse{resp}, err
}

func NewImageManager() Manager {
	return &imageManager{}
}
