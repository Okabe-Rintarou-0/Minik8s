package podworker

import (
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/kubelet/src/runtime/container"
)

type podWorkType byte

const (
	none podWorkType = iota
	podCreate
	podDelete
	podContainerStart
	podContainerCreateAndStart
	podContainerRemove
	podContainerRestart
)

type podWorkArgs interface{}

type podWork struct {
	WorkType podWorkType
	Arg      podWorkArgs
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

func newPodContainerStartWork(podUID types.UID, ID container.ID) podWork {
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

func newPodContainerRemoveWork(podUID types.UID, ID container.ID) podWork {
	return podWork{
		WorkType: podContainerRemove,
		Arg:      podContainerRemoveFnArg{podUID, ID},
	}
}

func newPodContainerRestartWork(pod *apiObject.Pod, ID container.ID, fullName string) podWork {
	return podWork{
		WorkType: podContainerRestart,
		Arg:      podContainerRestartFnArg{pod, ID, fullName},
	}
}
