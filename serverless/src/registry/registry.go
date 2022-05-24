// Package src docker registry
package registry

import (
	"context"
	"github.com/docker/distribution"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"log"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/serverless/src/utils"
	"os"
)

const (
	mediaType        = "application/vnd.docker.distribution.manifest.v2+json"
	registryImage    = "registry:2.8.0"
	registryName     = "local-registry"
	registryHost     = "127.0.0.1:5000"
	registryHostIP   = "127.0.0.1"
	registryHostPort = "5000"
)

var (
	cli *client.Client
)

func InitRegistry(host string) {
	utils.PullImg(registryImage)
	err := distribution.RegisterManifestSchema(mediaType, nil)
	if err != nil {
		log.Print(err)
		return
	}

	id, state := utils.FindContainer(registryName)

	if id == "" {
		id = utils.CreateContainer(registryName, &container.ContainerCreateConfig{
			Image:      registryImage,
			Entrypoint: nil,
			Cmd:        nil,
			Env:        nil,
			Volumes:    nil,
			Labels:     nil,
			IpcMode:    "",
			PidMode:    "",
			ExposedPorts: nat.PortSet{
				registryHostPort + "/tcp": {},
			},
			Tty:         false,
			Links:       nil,
			NetworkMode: "",
			Binds:       nil,
			PortBindings: nat.PortMap{
				registryHostPort + "/tcp": []nat.PortBinding{
					{
						HostIP:   registryHostIP,
						HostPort: registryHostPort,
					},
				},
			},
			VolumesFrom: nil,
		})
	}

	if state != "running" {
		utils.StartContainer(id)
	}

	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err == nil {
		log.Printf("init registry complete")
	} else {
		log.Print(err)
	}
}

func PushImage(image string) {
	pushReader, err := cli.ImagePush(context.Background(), image, types.ImagePushOptions{
		All:           false,
		RegistryAuth:  "",
		PrivilegeFunc: nil,
	})
	if err != nil {
		log.Printf("push image %s to registry error: %s\n", image, err)
	}
	wr, err := io.Copy(os.Stdout, pushReader)
	log.Print(wr)
}
