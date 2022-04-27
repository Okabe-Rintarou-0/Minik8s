package pod

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	"strconv"
	"testDocker/apiObject"
	"testDocker/kubelet/src/runtime/container"
	"testDocker/kubelet/src/runtime/image"
	"testDocker/kubelet/src/types"
	"time"
)

type Pod struct {
	ID         types.UID
	Name       string
	Namespace  string
	Containers []*container.Container
}

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
	GetPodStatus(pod *apiObject.Pod) (*PodStatus, error)
}

type podManager struct {
	cm container.ContainerManager
	im image.ImageManager
}

// needPullImage judges whether we need pull the image of given container spec
func (pm *podManager) needPullImage(container *apiObject.Container) (bool, error) {
	if container.ImagePullPolicy == apiObject.PullPolicyAlways {
		return true, nil
	}
	exist, err := pm.im.ExistsImage(container.Image)
	return !exist, err
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
func (pm *podManager) toVolumeBinds() []string {
	/// TODO implement it
	return nil
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

func (pm *podManager) addPortBindings(portBindings container.PortBindings, ports []apiObject.ContainerPort) error {
	for _, port := range ports {
		if port.Protocol == "" {
			port.Protocol = "tcp"
		}
		containerPort, err := nat.NewPort(port.Protocol, port.ContainerPort)
		if err != nil {
			return err
		}
		if port.HostIP == "" {
			port.HostIP = "127.0.0.1"
		}
		portBindings[containerPort] = []nat.PortBinding{{
			HostIP:   port.HostIP,
			HostPort: port.HostPort,
		}}
	}
	return nil
}

func (pm *podManager) addPortSet(portSet container.PortSet, ports []apiObject.ContainerPort) {
	for _, port := range ports {
		portSet[container.Port(port.ContainerPort+"/tcp")] = struct{}{}
	}
}

func (pm *podManager) getPauseContainerCreateConfig(pod *apiObject.Pod) (*container.ContainerCreateConfig, error) {
	labels := map[string]string{
		KubernetesPodUIDLabel: pod.UID(),
	}

	portBindings := container.PortBindings{}
	portSet := container.PortSet{}
	for _, c := range pod.Spec.Containers {
		err := pm.addPortBindings(portBindings, c.Ports)
		if err != nil {
			return nil, err
		}
		pm.addPortSet(portSet, c.Ports)
	}

	return &container.ContainerCreateConfig{
		Image:        pauseImage,
		Volumes:      nil,
		Labels:       labels,
		Binds:        nil,
		IpcMode:      "shareable",
		ExposedPorts: portSet,
		PortBindings: portBindings,
	}, nil
}

func (pm *podManager) getCommonContainerCreateConfig(c *apiObject.Container, podFullName string, podUID types.UID) *container.ContainerCreateConfig {
	// the label of given podUID
	labels := map[string]string{
		KubernetesPodUIDLabel: podUID,
	}
	pauseContainerFullName := pm.pauseContainerFullName(podFullName, podUID)
	pauseContainerRef := pm.toPauseContainerReference(podFullName, podUID)
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
}

func (pm *podManager) inspectionToContainerStatus(inspection *container.ContainerInspectInfo) (*container.ContainerStatus, error) {
	state := container.ContainerStateUnknown
	switch inspection.State.Status {
	case "running":
		state = container.ContainerStateRunning
	case "created":
		state = container.ContainerStateCreated
	case "exited":
		state = container.ContainerStateExited
	}

	createdAt, err := time.Parse(time.RFC3339Nano, inspection.Created)
	if err != nil {
		return nil, err
	}

	startedAt, err := time.Parse(time.RFC3339Nano, inspection.State.StartedAt)
	if err != nil {
		return nil, err
	}

	finishedAt, err := time.Parse(time.RFC3339Nano, inspection.State.FinishedAt)
	if err != nil {
		return nil, err
	}

	return &container.ContainerStatus{
		ID:           inspection.ID,
		Name:         inspection.Name,
		State:        state,
		CreatedAt:    createdAt,
		StartedAt:    startedAt,
		FinishedAt:   finishedAt,
		ExitCode:     inspection.State.ExitCode,
		ImageID:      inspection.Image,
		RestartCount: inspection.RestartCount,
		Error:        inspection.State.Error,
	}, nil
}

func (pm *podManager) getPodContainerStatuses(pod *apiObject.Pod) ([]*container.ContainerStatus, error) {
	containers, err := pm.cm.ListContainers(&container.ContainerListConfig{
		All: true,
		LabelSelector: container.LabelSelector{
			KubernetesPodUIDLabel: pod.UID(),
		},
	})
	if err != nil {
		return nil, err
	}

	containerStatuses := make([]*container.ContainerStatus, len(containers))
	for i, c := range containers {
		inspection, err := pm.cm.InspectContainer(c.ID)
		if err != nil {
			return nil, err
		}
		containerStatuses[i], err = pm.inspectionToContainerStatus(&inspection)
		if err != nil {
			return nil, err
		}
	}
	return containerStatuses, nil
}

// startPauseContainer starts the pause container that other common containers need
func (pm *podManager) startPauseContainer(pod *apiObject.Pod) error {
	// Step 1: Do we need pull the image?
	exists, err := pm.im.ExistsImage(pauseImage)
	if err != nil {
		return err
	}

	// Step 2: If needed, pull the image for the given container
	if !exists {
		fmt.Println("Need to pull image", pauseImage)
		err = pm.im.PullImage(pauseImage, &image.ImagePullConfig{
			Verbose: true,
			All:     false,
		})
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("No need to pull image %s, continue\n", pauseImage)
	}

	// Prepare
	podFullName := pod.FullName()
	podUID := pod.UID()

	// Step 3: Create a container
	fmt.Println("Now create the container")

	containerFullName := pm.pauseContainerFullName(podFullName, podUID)

	// get the container create config of pause
	var createConfig *container.ContainerCreateConfig
	createConfig, err = pm.getPauseContainerCreateConfig(pod)
	if err != nil {
		return err
	}

	var ID container.ContainerID
	ID, err = pm.cm.CreateContainer(containerFullName, createConfig)
	if err != nil {
		return err
	}
	fmt.Println("Create the container successfully, got ID", ID)

	// Step 4: Start this container
	fmt.Println("Now start the container with ID", ID)
	err = pm.cm.StartContainer(ID, &container.ContainerStartConfig{})
	return err
}

// startCommonContainer starts a common container according to the given spec
func (pm *podManager) startCommonContainer(pod *apiObject.Pod, c *apiObject.Container) error {
	// Step 1: Do we need pull the image?
	needPull, err := pm.needPullImage(c)
	if err != nil {
		return err
	}

	// Step 2: If needed, pull the image for the given container
	if needPull {
		fmt.Println("Need to pull image", c.Image)
		err = pm.im.PullImage(c.Image, &image.ImagePullConfig{
			Verbose: true,
			All:     false,
		})
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("No need to pull image %s, continue\n", c.Image)
	}

	// Prepare
	podFullName := pod.FullName()
	podUID := pod.UID()

	// Step 3: Create a container
	fmt.Println("Now create the container")

	containerFullName := pm.containerFullName(c.Name, podFullName, podUID, 0)

	var ID container.ContainerID
	ID, err = pm.cm.CreateContainer(containerFullName, pm.getCommonContainerCreateConfig(c, podFullName, podUID))
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

// CreatePod create a pod according to the given api object
func (pm *podManager) CreatePod(pod *apiObject.Pod) error {
	// Step 1: Start pause container
	err := pm.startPauseContainer(pod)
	if err != nil {
		return err
	}

	// Step 2: Start init containers
	/// TODO implement it

	// Step 3: Start common containers
	for _, c := range pod.Spec.Containers {
		err = pm.startCommonContainer(pod, &c)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Pod with UID %s created!\n", pod.UID())
	return nil
}

func (pm *podManager) GetPodStatus(pod *apiObject.Pod) (*PodStatus, error) {
	containerStatuses, err := pm.getPodContainerStatuses(pod)
	if err != nil {
		return nil, err
	}
	return &PodStatus{
		ID:                pod.UID(),
		Name:              pod.Name(),
		Namespace:         pod.Namespace(),
		IPs:               nil,
		ContainerStatuses: containerStatuses,
	}, nil
}
