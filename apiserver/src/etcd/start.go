package etcd

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/go-connections/nat"
	"log"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/kubelet/src/runtime/docker"
	"minik8s/kubelet/src/runtime/image"
)

const etcdImage = "bitnami/etcd:3.5"
const serverName = "etcd-server"

var (
	cm = container.NewContainerManager()
	im = image.NewImageManager()
)

func pullImage() {
	if exists, err := im.ExistsImage(etcdImage); !exists && err == nil {
		err := im.PullImage(etcdImage, &image.ImagePullConfig{All: false, Verbose: true})
		if err != nil {
			log.Print(err.Error())
			return
		}
	} else if err == nil {
		log.Print("[etcd] Image exists")
	} else {
		log.Print(err.Error())
		return
	}
}

func findContainer() (string, string) {
	ops := types.ContainerListOptions{All: true}
	ops.Filters = filters.NewArgs()
	ops.Filters.Add("name", serverName)
	containers, err := docker.Client.ContainerList(docker.Ctx, ops)
	if containers != nil && len(containers) > 0 && err == nil {
		return containers[0].ID, containers[0].State
	}
	return "", ""
}

func createContainer() string {
	ID, err := cm.CreateContainer(serverName, &container.ContainerCreateConfig{
		Image:      etcdImage,
		Entrypoint: nil,
		Cmd:        nil,
		Env:        []string{"ALLOW_NONE_AUTHENTICATION=yes", "ETCD_ADVERTISE_CLIENT_URLS=http://etcd-server:2379"},
		Volumes:    nil,
		ExposedPorts: nat.PortSet{
			"2379/tcp": {},
		},
		Tty:         false,
		Links:       nil,
		NetworkMode: "",
		Binds:       nil,
		PortBindings: nat.PortMap{
			"2379/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "2379",
				},
			},
		},
		VolumesFrom: nil,
	})
	if err != nil {
		log.Print(err.Error())
		return ""
	}
	return ID
}

func startContainer(ID string) {
	err := cm.StartContainer(ID, &container.ContainerStartConfig{})
	if err != nil {
		log.Print(err.Error())
		return
	}
}

func Start() {
	pullImage()
	ID, State := findContainer()
	if ID != "" {
		log.Printf("[etcd] Find etcd container, ID: %s, State: %s\n", ID, State)

	} else {
		ID = createContainer()
		log.Printf("[etcd] Create Container, ID: %s\n", ID)
	}
	if State != "running" && ID != "" {
		startContainer(ID)
		log.Printf("[etcd] Container starts\n")
	}
}
