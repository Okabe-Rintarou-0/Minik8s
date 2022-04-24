package container

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testDocker/src/runtime/image"
	"testing"
	"time"
)

func TestListContainer(t *testing.T) {
	cm := NewContainerManager()
	containers, err := cm.ListContainer(&ContainerListConfig{
		All: true,
	})
	assert.Nil(t, err)
	for _, c := range containers {
		fmt.Println(c.Name, c.ID, c.Image)
	}
}

func TestCreateStartAndRemoveContainer(t *testing.T) {
	cm := NewContainerManager()
	is := image.NewImageService()
	testImage := "nginx:latest"
	if exists, err := is.ExistsImage(testImage); !exists && err == nil {
		fmt.Printf("Image %s does not exist, so try to pull it\n", testImage)
		assert.Nil(t, is.PullImage(testImage, &image.ImagePullConfig{
			Verbose: true,
			All:     false,
		}))
	} else if err == nil {
		fmt.Printf("Image %s exists, continue\n", testImage)
	} else {
		fmt.Println(err.Error())
	}
	ID, err := cm.CreateContainer("test", &ContainerCreateConfig{
		Image:        testImage,
		Entrypoint:   nil,
		Cmd:          nil,
		Env:          nil,
		Volumes:      nil,
		ExposedPorts: nil,
		Tty:          false,
		Links:        nil,
		NetworkMode:  "",
		Binds:        nil,
		PortBindings: nil,
		VolumesFrom:  nil,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
	fmt.Printf("Create a container with ID %s\n", ID)

	fmt.Printf("Now inspect the container with ID %s\n", ID)
	inspection, err := cm.InspectContainer(ID)
	fmt.Println("Got inspection", inspection.Name, inspection.ID, inspection.Image)
	assert.Nil(t, err)

	fmt.Printf("Now start the container with ID %s\n", ID)
	err = cm.StartContainer(ID, &ContainerStartConfig{})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)

	fmt.Printf("Now stop the container with ID %s\n", ID)
	err = cm.StopContainer(ID, &ContainerStopConfig{timeout: time.Second})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)

	fmt.Printf("Now remove the container with ID %s\n", ID)
	err = cm.RemoveContainer(ID, &ContainerRemoveConfig{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
}
