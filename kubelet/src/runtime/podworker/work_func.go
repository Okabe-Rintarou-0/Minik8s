package podworker

import (
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/kubelet/src/runtime/container"
)

//CreatePod(pod *apiObject.Pod) error
//PodRemoveContainer(podUID types.UID, ID container.ID) error
//PodCreateAndStartContainer(pod *apiObject.Pod, target *apiObject.Container) error
//PodStartContainer(podUID types.UID, ID container.ID) error
//PodRestartContainer(pod *apiObject.Pod, containerID container.ID, fullName string) error

type PodCreateFn func(pod *apiObject.Pod) error
type PodDeleteFn func(pod *apiObject.Pod) error
type PodContainerStartFn func(podUID types.UID, ID container.ID) error
type PodContainerCreateAndStartFn func(pod *apiObject.Pod, target *apiObject.Container) error
type PodContainerRemoveFn func(podUID types.UID, ID container.ID) error
type PodContainerRestartFn func(pod *apiObject.Pod, containerID container.ID, fullName string) error

type podCreateFnArg struct {
	pod *apiObject.Pod
}

type podDeleteFnArg struct {
	pod *apiObject.Pod
}

type podContainerStartFnArg struct {
	podUID types.UID
	ID     container.ID
}

type podContainerCreateAndStartFnArg struct {
	pod    *apiObject.Pod
	target *apiObject.Container
}

type podContainerRemoveFnArg struct {
	podUID types.UID
	ID     container.ID
}

type podContainerRestartFnArg struct {
	pod      *apiObject.Pod
	ID       container.ID
	fullName string
}
