// Package src docker registry
// Ref: https://github.com/distribution/distribution/blob/main/docs/spec/api.md#deleting-an-image
// Ref: https://stackoverflow.com/questions/25436742/how-to-delete-images-from-a-private-docker-registry
package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/serverless/src/utils"
	"net/http"
	"os"

	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	rclient "github.com/docker/distribution/registry/client"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	RegistryImage = "registry:2.8.0"
	RegistryName  = "local-registry"
	RegistryHost     = "10.119.11.101:5000"
	//RegistryHost     = "0.0.0.0:5000"
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
	fmt.Println("init registry complete")
}

func PushImage(image string) error {
	fmt.Printf("Now push image %s\n", image)
	authConfig := types.AuthConfig{Username: "docker", Password: ""}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	pushReader, err := cli.ImagePush(context.Background(), image, types.ImagePushOptions{
		All:           false,
		RegistryAuth:  authStr,
		PrivilegeFunc: nil,
	})
	if err != nil {
		fmt.Printf("push image %s to registry error: %s\n", image, err)
		return err
	}
	wr, err := io.Copy(os.Stdout, pushReader)
	fmt.Println(wr)
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

// force (bool): Force removal of the image
// noprune (bool): Do not delete untagged parents
func RemoveImage(image string) error {
	// _, err := client.ImageRemove(context.Background(), "image_id", types.ImageRemoveOptions{})
	// if !errdefs.IsSystem(err) {
	// 	t.Fatalf("expected a Server Error, got %[1]T: %[1]v", err)
	// }
	imageTag := "latest"
	repo, err := newRepository(image)
	fmt.Print("check:",repo.Named())
	fmt.Println("newRepository completed")
	if err != nil {
		return err
	}

	ctx := context.Background()
	fmt.Println("chpt0")
	tagService := repo.Tags(ctx)
	fmt.Println("chpt1")
	desc, err := tagService.Get(ctx, imageTag)
	fmt.Println()
	if err != nil {
		return err
	}
	fmt.Println("chpt2")
	manifestService, err := repo.Manifests(ctx, nil)
	fmt.Println("chpt3")
	if err != nil {
		return err
	}

	return manifestService.Delete(ctx, desc.Digest)
}

func newRepository(imageName string) (distribution.Repository, error) {
	ref, err := reference.Parse(imageName)
	if err != nil {
		return nil, err
	}
	return rclient.NewRepository(ref.(reference.Named), "http://"+RegistryHost, http.DefaultTransport)

}
