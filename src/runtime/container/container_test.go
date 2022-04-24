package container

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestCreateAndRemoveContainer(t *testing.T) {
	cm := NewContainerManager()
	ID, err := cm.CreateContainer("test", &ContainerCreateConfig{
		Image:        "nginx:latest",
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
