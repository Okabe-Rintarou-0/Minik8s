// Package src docker registry
package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/serverless/src/utils"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	RegistryImage    = "registry:2.8.0"
	RegistryName     = "local-registry"
	RegistryHost     = "0.0.0.0:5000"
	RegistryHostIP   = "0.0.0.0"
	RegistryHostPort = "5000"
)

var (
	cli *client.Client
)

func InitRegistry() {
	utils.PullImg(RegistryImage)

	id, state, _ := utils.FindContainer(RegistryName)

	if id == "" {
		id = utils.CreateContainer(RegistryName, &container.ContainerCreateConfig{
			Image:      RegistryImage,
			Entrypoint: nil,
			Cmd:        nil,
			Env:        nil,
			Volumes:    nil,
			Labels:     nil,
			IpcMode:    "",
			PidMode:    "",
			ExposedPorts: nat.PortSet{
				RegistryHostPort + "/tcp": {},
			},
			Tty:         false,
			Links:       nil,
			NetworkMode: "",
			Binds:       nil,
			PortBindings: nat.PortMap{
				RegistryHostPort + "/tcp": []nat.PortBinding{
					{
						HostIP:   RegistryHostIP,
						HostPort: RegistryHostPort,
					},
				},
			},
			VolumesFrom: nil,
		})
	}

	if state != "running" {
		utils.StartContainer(id)
	}

	cli, _ = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	log.Printf("init registry complete")
}

func PushImage(image string) error {
	authConfig := types.AuthConfig{Username: "docker", Password: ""}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		log.Print(err)
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	pushReader, err := cli.ImagePush(context.Background(), image, types.ImagePushOptions{
		All:           false,
		RegistryAuth:  authStr,
		PrivilegeFunc: nil,
	})
	if err != nil {
		log.Printf("push image %s to registry error: %s\n", image, err)
		return err
	}
	wr, err := io.Copy(os.Stdout, pushReader)
	log.Print(wr)
	return nil
}

func PullImage(image string) {
	authConfig := types.AuthConfig{Username: "docker", Password: ""}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		log.Print(err)
		return
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	pullReader, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{
		All:           false,
		RegistryAuth:  authStr,
		PrivilegeFunc: nil,
	})
	wr, err := io.Copy(os.Stdout, pullReader)
	log.Print(wr)
}
