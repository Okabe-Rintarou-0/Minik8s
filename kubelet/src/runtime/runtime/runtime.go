package runtime

import (
	"fmt"
	"minik8s/apiObject"
	"minik8s/kubelet/src/podutil"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/kubelet/src/runtime/image"
	"minik8s/kubelet/src/types"
	"strconv"
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

func (podStatus *PodStatus) GetContainerStatusByName(name string) *container.ContainerStatus {
	for _, cs := range podStatus.ContainerStatuses {
		if cs.Name == name {
			return cs
		}
	}
	return nil
}

type Pods []Pod

type PodStatuses = map[types.UID]*PodStatus

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

type Manager interface {
	CreatePod(pod *apiObject.Pod) error
	GetPodStatus(pod *apiObject.Pod) (*PodStatus, error)
	GetPodStatuses() (PodStatuses, error)
	PodRemoveContainer(podUID types.UID, ID container.ContainerID) error
	PodCreateAndStartContainer(pod *apiObject.Pod, target *apiObject.Container) error
	PodStartContainer(podUID types.UID, ID container.ContainerID) error
	PodRestartContainer(pod *apiObject.Pod, containerID container.ContainerID, fullName string) error
}

type runtimeManager struct {
	cm container.Manager
	im image.ImageManager
}

func (rm *runtimeManager) PodCreateAndStartContainer(pod *apiObject.Pod, target *apiObject.Container) error {
	return rm.startCommonContainer(pod, target)
}

func (rm *runtimeManager) PodStartContainer(podUID types.UID, ID container.ContainerID) error {
	return rm.cm.StartContainer(ID, &container.ContainerStartConfig{})
}

func (rm *runtimeManager) PodRestartContainer(pod *apiObject.Pod, containerID container.ContainerID, fullName string) error {
	parseSucc, _, _, _, _, restartCount := podutil.ParseContainerFullName(fullName)
	if !parseSucc {
		panic("Could not happen")
	}
	newName := fullName[:len(fullName)-1]
	newName += strconv.Itoa(restartCount + 1)
	err := rm.cm.RenameContainer(containerID, newName)
	if err != nil {
		return err
	}

	return rm.cm.StartContainer(containerID, &container.ContainerStartConfig{})
}

func (rm *runtimeManager) PodRemoveContainer(podUID types.UID, ID container.ContainerID) error {
	return rm.cm.RemoveContainer(ID, &container.ContainerRemoveConfig{})
}

// CreatePod create a pod according to the given api object
func (rm *runtimeManager) CreatePod(pod *apiObject.Pod) error {
	// Step 1: Start pause container
	err := rm.startPauseContainer(pod)
	if err != nil {
		return err
	}

	// Step 2: Start init containers
	/// TODO implement it

	// Step 3: Start common containers
	for _, c := range pod.Spec.Containers {
		err = rm.startCommonContainer(pod, &c)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Pod with UID %s created!\n", pod.UID())
	return nil
}

func (rm *runtimeManager) GetPodStatus(pod *apiObject.Pod) (*PodStatus, error) {
	containerStatuses, err := rm.getPodContainerStatuses(pod)
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

func (rm *runtimeManager) GetPodStatuses() (PodStatuses, error) {
	allContainerStatuses, err := rm.getAllPodContainers()
	if err != nil {
		return nil, err
	}
	podStatuses := make(PodStatuses)
	for podUID, cs := range allContainerStatuses {
		podStatuses[podUID] = &PodStatus{
			ID:                podUID,
			ContainerStatuses: cs,
		}
	}
	return podStatuses, nil
}

func NewPodManager() Manager {
	return &runtimeManager{
		cm: container.NewContainerManager(),
		im: image.NewImageManager(),
	}
}
