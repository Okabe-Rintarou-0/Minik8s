package image

import (
	"bufio"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"io"
	"testDocker/src/runtime/docker"
)

type ImageSummary = types.ImageSummary

type ImageService interface {
	PullImage(name string, config *ImagePullConfig) error
	ExistsImage(name string) (bool, error)
	ListImages(config *ImageListConfig) ([]ImageSummary, error)
	ListImagesByName(name string, config *ImageListConfig) ([]ImageSummary, error)
	RemoveImage(ID string, config *ImageRemoveConfig) (ImageRemoveResponse, error)
}

type imageService struct {
}

func (is *imageService) ExistsImage(name string) (bool, error) {
	images, err := is.ListImagesByName(name, &ImageListConfig{All: true})
	return len(images) > 0, err
}

func (is *imageService) ListImages(config *ImageListConfig) ([]ImageSummary, error) {
	return docker.Client.ImageList(docker.Ctx, types.ImageListOptions{
		All: config.All,
	})
}

func (is *imageService) ListImagesByName(name string, config *ImageListConfig) ([]ImageSummary, error) {
	filter := filters.NewArgs()
	filter.Add("dangling", "false")
	filter.Add("reference", name)
	return docker.Client.ImageList(docker.Ctx, types.ImageListOptions{
		All:     config.All,
		Filters: filter,
	})
}

func (is *imageService) PullImage(name string, config *ImagePullConfig) error {
	resp, err := docker.Client.ImagePull(docker.Ctx, name, types.ImagePullOptions{
		All: config.All,
	})
	if err != nil {
		return err
	}
	r := bufio.NewReader(resp)
	if config.Verbose {
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

func (is *imageService) RemoveImage(ID string, config *ImageRemoveConfig) (ImageRemoveResponse, error) {
	resp, err := docker.Client.ImageRemove(docker.Ctx, ID, *config)
	return ImageRemoveResponse{resp}, err
}

func NewImageService() ImageService {
	return &imageService{}
}
