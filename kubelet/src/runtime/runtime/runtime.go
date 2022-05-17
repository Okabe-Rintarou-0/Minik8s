package runtime

import (
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/entity"
	"minik8s/kubelet/src/podutil"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/kubelet/src/runtime/image"
	"minik8s/util/logger"
	"os"
	"strconv"
	"time"
)

var log = logger.Log("Runtime")

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
	// PodLifecycle of containers in the pod.
	ContainerStatuses []*container.Status
}

func (podStatus *PodStatus) ToEntity() *entity.PodStatus {
	hostname, _ := os.Hostname()
	cpuPercent, memPercent := calcMetrics(podStatus.ContainerStatuses)
	return &entity.PodStatus{
		ID:         podStatus.ID,
		Name:       podStatus.Name,
		Node:       hostname,
		Namespace:  podStatus.Namespace,
		Lifecycle:  entity.PodRunning,
		CpuPercent: cpuPercent,
		MemPercent: memPercent,
		Error:      "",
		SyncTime:   time.Now(),
	}
}

func (podStatus *PodStatus) GetContainerStatusByName(name string) *container.Status {
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
func (pod *Pod) GetContainerByID(ID container.ID) *container.Container {
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
	DeletePod(pod *apiObject.Pod) error
	GetPodStatus(pod *apiObject.Pod) (*PodStatus, error)
	GetPodStatuses() (PodStatuses, error)
	PodRemoveContainer(podUID types.UID, ID container.ID) error
	PodCreateAndStartContainer(pod *apiObject.Pod, target *apiObject.Container) error
	PodStartContainer(podUID types.UID, ID container.ID) error
	PodRestartContainer(pod *apiObject.Pod, containerID container.ID, fullName string) error
}

type runtimeManager struct {
	cm container.Manager
	im image.Manager
}

func (rm *runtimeManager) PodCreateAndStartContainer(pod *apiObject.Pod, target *apiObject.Container) error {
	return rm.startCommonContainer(pod, target)
}

func (rm *runtimeManager) PodStartContainer(podUID types.UID, ID container.ID) error {
	return rm.cm.StartContainer(ID, &container.StartConfig{})
}

func (rm *runtimeManager) PodRestartContainer(pod *apiObject.Pod, containerID container.ID, fullName string) error {
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

	return rm.cm.StartContainer(containerID, &container.StartConfig{})
}

func (rm *runtimeManager) PodRemoveContainer(podUID types.UID, ID container.ID) error {
	return rm.cm.RemoveContainer(ID, &container.RemoveConfig{})
}

// DeletePod deletes a pod according to the given api object
func (rm *runtimeManager) DeletePod(pod *apiObject.Pod) error {
	log("Delete pod[ID = %s]", pod.UID())
	// Step 1: Remove common container
	err := rm.removePodCommonContainers(pod)

	if err != nil {
		return err
	}

	// Step 2: Remove init containers
	/// TODO implement it

	// Step 3: Remove pause containers
	err = rm.removePauseContainer(pod)
	if err != nil {
		return err
	}

	log("Pod[ID = %s] has been removed!", pod.UID())
	return nil
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

	log("Pod[ID = %s] has been created!", pod.UID())
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
	allContainerStatuses := rm.getAllPodContainers()
	podStatuses := make(PodStatuses)
	for podUID, cs := range allContainerStatuses {
		podStatuses[podUID] = &PodStatus{
			ID:                podUID,
			ContainerStatuses: cs,
		}
		//fmt.Printf("Convert to entity would be: %v\n", *podStatuses[podUID].ToEntity())
	}
	return podStatuses, nil
}

func NewPodManager() Manager {
	return &runtimeManager{
		cm: container.NewContainerManager(),
		im: image.NewImageManager(),
	}
}
