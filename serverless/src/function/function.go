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
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	pythonImage     = "python:3.10-slim"
	exposedPort     = "8080"
	dockerfilePath  = "serverless/src/app/Dockerfile"
	requirementPath = "serverless/src/app/requirements.txt"
	mainCodePath    = "serverless/src/app/main.py" // a fixed template, start a http server
	dockerfile      = "Dockerfile"
	funcCode        = "func.py" // rename to "func.py" in docker
	mainCode        = "main.py"
	requirement     = "requirements.txt"
)

func CreateFunctionImage(name string, codePath string) error {
	err := utils.PullImg(pythonImage)
	if err != nil {
		return err
	}
	err = createImage(name, codePath)
	if err != nil {
		return err
	}
	imageName := registry.RegistryHost + "/" + name
	return registry.PushImage(imageName)
	//_, _ = createContainer(name, containerName, imageName)
}

func RemoveFunctionImage(name string) error {
	//TODO implement it
	return nil
}

func createContainer(name, containerName, imageName string) (string, string) {
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
					HostIP: registry.RegistryHostIP,
				},
			},
		},
		VolumesFrom: nil,
	})
	utils.StartContainer(id)
	id, state, port := utils.FindContainer(containerName)
	portStr := strconv.FormatUint(uint64(port), 10)
	log.Printf("Ready to serve %s, container id: %s, container state: %s, host port: %s\n", name, id, state, portStr)
	return id, portStr
}

func copyFile(tw *tar.Writer, path string, filename string) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	readFile, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	tarHeader := &tar.Header{
		Name: filename,
		Size: int64(len(readFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return err
	}
	_, err = tw.Write(readFile)
	if err != nil {
		return err
	}
	return nil
}

func createImage(name, funcCodePath string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	if err = copyFile(tw, dockerfilePath, dockerfile); err != nil {
		return err
	}
	if err = copyFile(tw, mainCodePath, mainCode); err != nil {
		return err
	}
	if err = copyFile(tw, funcCodePath, funcCode); err != nil {
		return err
	}
	if err = copyFile(tw, requirementPath, requirement); err != nil {
		return err
	}

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
		return err
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		return err
	}
	return nil
}
