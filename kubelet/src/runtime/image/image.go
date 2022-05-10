package image

import (
	"bufio"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"io"
	"minik8s/kubelet/src/runtime/docker"
)

type Summary = types.ImageSummary

type Manager interface {
	PullImage(name string, config *ImagePullConfig) error
	ExistsImage(name string) (bool, error)
	ListImages(config *ImageListConfig) ([]Summary, error)
	ListImagesByName(name string, config *ImageListConfig) ([]Summary, error)
	RemoveImage(ID string, config *ImageRemoveConfig) (ImageRemoveResponse, error)
}

type imageManager struct {
}

func (is *imageManager) ExistsImage(name string) (bool, error) {
	images, err := is.ListImagesByName(name, &ImageListConfig{All: true})
	return len(images) > 0, err
}

func (is *imageManager) ListImages(config *ImageListConfig) ([]Summary, error) {
	return docker.Client.ImageList(docker.Ctx, types.ImageListOptions{
		All: config.All,
	})
}

func (is *imageManager) ListImagesByName(name string, config *ImageListConfig) ([]Summary, error) {
	filter := filters.NewArgs()
	filter.Add("dangling", "false")
	filter.Add("reference", name)
	return docker.Client.ImageList(docker.Ctx, types.ImageListOptions{
		All:     config.All,
		Filters: filter,
	})
}

func (is *imageManager) PullImage(name string, config *ImagePullConfig) error {
	resp, err := docker.Client.ImagePull(docker.Ctx, name, types.ImagePullOptions{
		All: config.All,
	})
	if err != nil {
		return err
	}
	defer resp.Close()
	if config.Verbose {
		r := bufio.NewReader(resp)
		for {
			row, err := r.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			fmt.Print(row)
		}
	}
	return nil
}

func (is *imageManager) RemoveImage(ID string, config *ImageRemoveConfig) (ImageRemoveResponse, error) {
	resp, err := docker.Client.ImageRemove(docker.Ctx, ID, *config)
	return ImageRemoveResponse{resp}, err
}

func NewImageManager() Manager {
	return &imageManager{}
}
