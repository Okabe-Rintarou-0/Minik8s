package image

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"io/ioutil"
	"testDocker/src/runtime/docker"
)

type ImageSummary = types.ImageSummary

type ImageService interface {
	PullImage(name string, config *ImagePullConfig) (response string, err error)
	ListImages(config *ImageListConfig) ([]ImageSummary, error)
	ListImagesByName(name string, config *ImageListConfig) ([]ImageSummary, error)
}

type imageService struct {
}

func NewImageService() ImageService {
	return &imageService{}
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

func (is *imageService) PullImage(name string, config *ImagePullConfig) (response string, err error) {
	resp, err := docker.Client.ImagePull(docker.Ctx, name, types.ImagePullOptions{
		All: config.All,
	})
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp)
	return string(body), err
}
