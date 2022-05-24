package function

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/serverless/src/registry"
	"minik8s/serverless/src/utils"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	uuid "github.com/satori/go.uuid"
)

const (
	pythonImage     = "python:3.10-slim"
	exposedPort     = "8080"
	dockerfilePath  = "../src/app/Dockerfile"
	codePath        = "../src/app/helloworld.py"
	requirementPath = "../src/app/requirements.txt"
	name            = "test1"
	dockerfile      = "Dockerfile"
	code            = "func.py"
	requirement     = "requirements.txt"
)

func InitFunction() {
	utils.PullImg(pythonImage)
	createImage()

	uid := uuid.NewV4().String()
	containerName := name + "-" + uid
	imageName := registry.RegistryHost + "/" + name

	registry.PushImage(imageName)
	createContainer(containerName, imageName)

}

func createContainer(containerName, imageName string) {
	id := utils.CreateContainer(containerName, &container.ContainerCreateConfig{
		Image:      imageName,
		Entrypoint: nil,
		Cmd:        nil,
		Env:        nil,
		Volumes:    nil,
		Labels:     nil,
		IpcMode:    "",
		PidMode:    "",
		ExposedPorts: nat.PortSet{
			exposedPort + "/tcp": {},
		},
		Tty:         false,
		Links:       nil,
		NetworkMode: "",
		Binds:       nil,
		PortBindings: nat.PortMap{
			exposedPort + "/tcp": []nat.PortBinding{
				{
					HostIP:   registry.RegistryHostIP,
					HostPort: "10303",
				},
			},
		},
		VolumesFrom: nil,
	})
	log.Printf("container created, id=%s\n", id)
	utils.StartContainer(id)
	log.Printf("Ready to serve " + name)
}

func copyFile(tw *tar.Writer, path string, filename string) {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err, " :unable to open file "+path)
	}
	readFile, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err, " :unable to read file "+path)
	}

	tarHeader := &tar.Header{
		Name: filename,
		Size: int64(len(readFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		log.Fatal(err, " :unable to write tar header")
	}
	_, err = tw.Write(readFile)
	if err != nil {
		log.Fatal(err, " :unable to write tar body")
	}
}

func createImage() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err, " :unable to init client")
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	copyFile(tw, dockerfilePath, dockerfile)
	copyFile(tw, codePath, code)
	copyFile(tw, requirementPath, requirement)

	tarReader := bytes.NewReader(buf.Bytes())

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		tarReader,
		types.ImageBuildOptions{
			Context:    tarReader,
			Dockerfile: dockerfile,
			Tags:       []string{registry.RegistryHost + "/" + name},
			Remove:     true})
	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
}
