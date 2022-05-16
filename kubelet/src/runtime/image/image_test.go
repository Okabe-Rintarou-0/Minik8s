package image

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImageService_ListImages(t *testing.T) {
	im := NewImageManager()
	images, err := im.ListImages(&ListConfig{All: true})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
	fmt.Println("Read readme-images: ")
	for _, image := range images {
		fmt.Println(image.ID)
	}
}

func TestImageService_ListImagesByName(t *testing.T) {
	im := NewImageManager()
	images, err := im.ListImagesByName("nginx:latest", &ListConfig{All: true})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
	fmt.Println("Read readme-images: ")
	for _, image := range images {
		fmt.Println(image.ID)
	}
}

func TestImageService_PullAndRemoveImage(t *testing.T) {
	im := NewImageManager()
	err := im.PullImage("nginx:latest", &PullConfig{All: false, Verbose: true})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)

	images, err := im.ListImagesByName("nginx:latest", &ListConfig{All: true})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
	fmt.Println("Now remove image with ID", images[0].ID)

	resp, err := im.RemoveImage(images[0].ID, &RemoveConfig{
		Force:         false,
		PruneChildren: false,
	})
	fmt.Println(resp.DeletedItems)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
}
