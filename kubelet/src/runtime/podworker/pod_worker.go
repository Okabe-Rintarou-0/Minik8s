package podworker

import (
	"fmt"
	"minik8s/apiObject"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/kubelet/src/types"
)

type podWorkType byte

const (
	podCreate podWorkType = iota
	podDelete
	podContainerStart
	podContainerCreateAndStart
	podContainerRemove
	podContainerRestart
)

//CreatePod(pod *apiObject.Pod) error
//PodRemoveContainer(podUID types.UID, ID container.ContainerID) error
//PodCreateAndStartContainer(pod *apiObject.Pod, target *apiObject.Container) error
//PodStartContainer(podUID types.UID, ID container.ContainerID) error
//PodRestartContainer(pod *apiObject.Pod, containerID container.ContainerID, fullName string) error

type PodCreateFn func(pod *apiObject.Pod) error
type PodDeleteFn func(pod *apiObject.Pod) error
type PodContainerStartFn func(podUID types.UID, ID container.ContainerID) error
type PodContainerCreateAndStartFn func(pod *apiObject.Pod, target *apiObject.Container) error
type PodContainerRemoveFn func(podUID types.UID, ID container.ContainerID) error
type PodContainerRestartFn func(pod *apiObject.Pod, containerID container.ContainerID, fullName string) error

type podCreateFnArg struct {
	pod *apiObject.Pod
}

type podDeleteFnArg struct {
	pod *apiObject.Pod
}

type podContainerStartFnArg struct {
	podUID types.UID
	ID     container.ContainerID
}

type podContainerCreateAndStartFnArg struct {
	pod    *apiObject.Pod
	target *apiObject.Container
}

type podContainerRemoveFnArg struct {
	podUID types.UID
	ID     container.ContainerID
}

type podContainerRestartFnArg struct {
	pod      *apiObject.Pod
	ID       container.ContainerID
	fullName string
}

type podWorkArgs interface{}

type podWork struct {
	WorkType podWorkType
	Arg      podWorkArgs
}

type podWorker struct {
	PodCreateFn                  PodCreateFn
	PodDeleteFn                  PodDeleteFn
	PodContainerStartFn          PodContainerStartFn
	PodContainerRestartFn        PodContainerRestartFn
	PodContainerCreateAndStartFn PodContainerCreateAndStartFn
	PodContainerRemoveFn         PodContainerRemoveFn
}

func newPodCreateWork(pod *apiObject.Pod) podWork {
	return podWork{
		WorkType: podCreate,
		Arg:      podCreateFnArg{pod},
	}
}

func newPodDeleteWork(pod *apiObject.Pod) podWork {
	return podWork{
		WorkType: podDelete,
		Arg:      podDeleteFnArg{pod},
	}
}

func newPodContainerStartWork(podUID types.UID, ID container.ContainerID) podWork {
	return podWork{
		WorkType: podContainerStart,
		Arg:      podContainerStartFnArg{podUID, ID},
	}
}

func newPodContainerCreateAndStartWork(pod *apiObject.Pod, target *apiObject.Container) podWork {
	return podWork{
		WorkType: podContainerCreateAndStart,
		Arg:      podContainerCreateAndStartFnArg{pod, target},
	}
}

func newPodContainerRemoveWork(podUID types.UID, ID container.ContainerID) podWork {
	return podWork{
		WorkType: podContainerRemove,
		Arg:      podContainerRemoveFnArg{podUID, ID},
	}
}

func newPodContainerRestartWork(pod *apiObject.Pod, ID container.ContainerID, fullName string) podWork {
	return podWork{
		WorkType: podContainerRestart,
		Arg:      podContainerRestartFnArg{pod, ID, fullName},
	}
}

func (w *podWorker) doWork(work podWork) {
	var err error
	switch work.WorkType {
	case podCreate:
		fmt.Println("pod worker received pod create job")
		arg := work.Arg.(podCreateFnArg)
		err = w.PodCreateFn(arg.pod)
	case podDelete:
		fmt.Println("pod worker received pod delete job")
		arg := work.Arg.(podDeleteFnArg)
		err = w.PodDeleteFn(arg.pod)
	case podContainerCreateAndStart:
		arg := work.Arg.(podContainerCreateAndStartFnArg)
		err = w.PodContainerCreateAndStartFn(arg.pod, arg.target)
	case podContainerRemove:
		arg := work.Arg.(podContainerRemoveFnArg)
		err = w.PodContainerRemoveFn(arg.podUID, arg.ID)
	case podContainerStart:
		arg := work.Arg.(podContainerStartFnArg)
		err = w.PodContainerStartFn(arg.podUID, arg.ID)
	case podContainerRestart:
		fmt.Println("pod worker received restart job")
		arg := work.Arg.(podContainerRestartFnArg)
		err = w.PodContainerRestartFn(arg.pod, arg.ID, arg.fullName)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (w *podWorker) Run(workCh <-chan podWork) {
	for {
		select {
		case work := <-workCh:
			w.doWork(work)
		}
	}
}
