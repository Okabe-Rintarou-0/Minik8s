package utils

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"log"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/kubelet/src/runtime/docker"
)

func FindContainer(containerName string) (string, string, uint16) {
	ops := types.ContainerListOptions{All: true}
	ops.Filters = filters.NewArgs()
	ops.Filters.Add("name", containerName)
	containers, err := docker.Client.ContainerList(docker.Ctx, ops)
	if containers != nil && len(containers) > 0 && err == nil {
		return containers[0].ID, containers[0].State, containers[0].Ports[0].PublicPort
	}
	if err != nil {
		log.Print(err.Error())
	}
	return "", "", 0
}

func CreateContainer(serverName string, config *container.ContainerCreateConfig) string {
	ID, err := cm.CreateContainer(serverName, config)
	if err != nil {
		log.Print(err.Error())
		return ""
	}
	return ID
}

func StartContainer(ID string) {
	err := cm.StartContainer(ID, &container.StartConfig{})
	if err != nil {
		log.Print(err.Error())
	}
}
