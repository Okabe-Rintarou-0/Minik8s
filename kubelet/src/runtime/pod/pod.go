package pod

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	"strconv"
	"testDocker/apiObject"
	"testDocker/kubelet/src/runtime/container"
	"testDocker/kubelet/src/runtime/image"
	"testDocker/kubelet/src/types"
)

type Pod struct {
	ID         types.UID
	Name       string
	Namespace  string
	Containers []*container.Container
}

type ContainerType = byte

const (
	PauseContainer ContainerType = iota
	CommonContainer
)

// PodStatus represents the status of the pod and its containers.
type PodStatus struct {
	// ID of the pod.
	ID types.UID
	// Name of the pod.
	Name string
	// Namespace of the pod.
	Namespace string
	// All IPs assigned to this pod
	IPs []string
	// Status of containers in the pod.
	ContainerStatuses []*container.ContainerStatus
}

type Pods []Pod

// FullName is the full name of the pod
func (pod *Pod) FullName() string {
	return pod.Name + "_" + pod.Namespace
}

// GetContainerByID returns the container of pod given the ID of it
func (pod *Pod) GetContainerByID(ID container.ContainerID) *container.Container {
	for _, c := range pod.Containers {
		if c.ID == ID {
			return c
		}
	}
	return nil
}

func (pod *Pod) GetContainerByName(name string) *container.Container {
	for _, c := range pod.Containers {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// GetPodByUID returns the Pod given the UID of it
func (pods Pods) GetPodByUID(ID types.UID) *Pod {
	for _, pod := range pods {
		if pod.ID == ID {
			return &pod
		}
	}
	return nil
}

func (pods Pods) GetPodByFullName(fullName string) *Pod {
	for _, pod := range pods {
		if pod.FullName() == fullName {
			return &pod
		}
	}
	return nil
}

type PodManager interface {
	CreatePod(pod *apiObject.Pod) error
}

type podManager struct {
	cm container.ContainerManager
	im image.ImageManager
}

func (pm *podManager) CreatePod(pod *apiObject.Pod) error {
	return nil
}

// needPullImage judges whether we need pull the image of given container spec
func (pm *podManager) needPullImage(container *apiObject.Container, containerType ContainerType) (bool, string, error) {
	var imageName string
	if containerType == PauseContainer {
		imageName = pauseImage
	} else {
		imageName = container.Image
	}
	if container != nil && container.ImagePullPolicy == apiObject.PullPolicyAlways {
		return true, imageName, nil
	}
	exist, err := pm.im.ExistsImage(imageName)
	return !exist, imageName, err
}

// toFormattedEnv changes containerEnv to adapted form, like "FOO=bar"
// where FOO is name and bar is value
func (pm *podManager) toFormattedEnv(containerEnv []apiObject.EnvVar) []string {
	var env []string
	for _, ev := range containerEnv {
		env = append(env, ev.Name+"="+ev.Value)
	}
	return env
}

// toVolumeBinds returns the binds of volumes
func (pm *podManager) toVolumeBinds(containerEnv []apiObject.EnvVar) []string {
	var env []string
	for _, ev := range containerEnv {
		env = append(env, ev.Name+"="+ev.Value)
	}
	return env
}

func (pm *podManager) containerFullName(containerName, podFullName string, podUID string, restartCount int) string {
	return "k8s_" + containerName + "_" + podFullName + "_" + podUID + "_" + strconv.Itoa(restartCount)
}

func (pm *podManager) pauseContainerFullName(podFullName string, podUID types.UID) string {
	return pm.containerFullName(pauseContainerName, podFullName, podUID, 0)
}

func (pm *podManager) toContainerReference(containerFullName string) string {
	return "container:" + containerFullName
}

func (pm *podManager) toPauseContainerReference(podFullName string, podUID types.UID) string {
	return pm.toContainerReference(pm.pauseContainerFullName(podFullName, podUID))
}

func (pm *podManager) toPortBindings(ports []apiObject.ContainerPort) container.PortBindings {
	portBindings := container.PortBindings{}
	for _, port := range ports {
		if port.Protocol == "" {
			port.Protocol = "TCP"
		}
		containerPort, err := nat.NewPort(port.Protocol, port.ContainerPort)
		if err != nil {
			return nil
		}
		portBindings[containerPort] = []nat.PortBinding{{
			HostIP:   port.HostIP,
			HostPort: port.HostPort,
		}}
	}
	return portBindings
}

func (pm *podManager) startCommonContainer(pod *apiObject.Pod, c *apiObject.Container) error {
	return pm.startContainer(c, pod, CommonContainer)
}

func (pm *podManager) startPauseContainer(pod *apiObject.Pod) error {
	return pm.startContainer(nil, pod, PauseContainer)
}

func (pm *podManager) getContainerCreateConfig(c *apiObject.Container, pod *apiObject.Pod, containerType ContainerType) *container.ContainerCreateConfig {
	// the label of given podUID
	labels := map[string]string{
		KubernetesPodUIDLabel: pod.UID(),
	}
	pauseContainerFullName := pm.pauseContainerFullName(pod.FullName(), pod.UID())
	pauseContainerRef := pm.toPauseContainerReference(pod.FullName(), pod.UID())

	switch containerType {
	case CommonContainer:
		return &container.ContainerCreateConfig{
			Image:       c.Image,
			Entrypoint:  c.Command,
			Cmd:         c.Args,
			Env:         pm.toFormattedEnv(c.Env),
			Volumes:     nil,
			Labels:      labels,
			Tty:         c.TTY,
			NetworkMode: container.NetworkMode(pauseContainerRef),
			IpcMode:     container.IpcMode(pauseContainerRef),
			PidMode:     container.PidMode(pauseContainerRef),
			Binds:       nil,
			VolumesFrom: []string{pauseContainerFullName},
		}
	case PauseContainer:
		return &container.ContainerCreateConfig{
			Image:   pauseImage,
			Volumes: nil,
			Labels:  labels,
			Binds:   nil,
		}
	}
	panic("invalid container type")
}

// startContainer starts a container according to the given spec
func (pm *podManager) startContainer(c *apiObject.Container, pod *apiObject.Pod, containerType ContainerType) error {
	// Step 1: Do we need pull the image?
	needPull, imageName, err := pm.needPullImage(c, containerType)
	if err != nil {
		return err
	}

	// Step 2: If needed, pull the image for the given container
	if needPull {
		fmt.Println("Need to pull image", imageName)
		err = pm.im.PullImage(imageName, &image.ImagePullConfig{
			Verbose: true,
			All:     false,
		})
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("No need to pull image %s, continue\n", imageName)
	}

	// Step 3: Create a container
	fmt.Println("Now create the container")
	var ID container.ContainerID
	var containerFullName string
	switch containerType {
	case CommonContainer:
		containerFullName = pm.containerFullName(c.Name, pod.FullName(), pod.UID(), 0)
	case PauseContainer:
		containerFullName = pm.pauseContainerFullName(pod.FullName(), pod.UID())
	default:
		panic("invalid container type")
	}
	ID, err = pm.cm.CreateContainer(containerFullName, pm.getContainerCreateConfig(c, pod, containerType))
	if err != nil {
		return err
	}
	fmt.Println("Create the container successfully, got ID", ID)

	// Step 4: Start this container
	fmt.Println("Now start the container with ID", ID)
	err = pm.cm.StartContainer(ID, &container.ContainerStartConfig{})
	return err
}

func NewPodManager() PodManager {
	return &podManager{
		cm: container.NewContainerManager(),
		im: image.NewImageManager(),
	}
}
