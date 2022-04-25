package image

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImageService_ListImages(t *testing.T) {
	is := NewImageService()
	images, err := is.ListImages(&ImageListConfig{All: true})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
	fmt.Println("Read images: ")
	for _, image := range images {
		fmt.Println(image.ID)
	}
}

func TestImageService_ListImagesByName(t *testing.T) {
	is := NewImageService()
	images, err := is.ListImagesByName("nginx:latest", &ImageListConfig{All: true})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
	fmt.Println("Read images: ")
	for _, image := range images {
		fmt.Println(image.ID)
	}
}

func TestImageService_PullAndRemoveImage(t *testing.T) {
	is := NewImageService()
	err := is.PullImage("nginx:latest", &ImagePullConfig{All: false, Verbose: true})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)

	images, err := is.ListImagesByName("nginx:latest", &ImageListConfig{All: true})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
	fmt.Println("Now remove image with ID", images[0].ID)

	resp, err := is.RemoveImage(images[0].ID, &ImageRemoveConfig{
		Force:         false,
		PruneChildren: false,
	})
	fmt.Println(resp.DeletedItems)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
}
