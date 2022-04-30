package podworker

import (
	"minik8s/apiObject"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/kubelet/src/runtime/pleg"
	"minik8s/kubelet/src/types"
)

const (
	workChanSize = 5
)

type workChanMap map[types.UID]chan podWork

type Manager interface {
	AddPod(pod *apiObject.Pod)
	DeletePod(pod *apiObject.Pod)
	UpdatePod(event *pleg.PodLifecycleEvent)
}

type manager struct {
	workChanMap                  workChanMap
	podCreateFn                  PodCreateFn
	podDeleteFn                  PodDeleteFn
	podContainerCreateAndStartFn PodContainerCreateAndStartFn
	podContainerStartFn          PodContainerStartFn
	podContainerRemoveFn         PodContainerRemoveFn
	podContainerRestartFn        PodContainerRestartFn
}

func (m *manager) newWorker() *podWorker {
	return &podWorker{
		PodCreateFn:                  m.podCreateFn,
		PodDeleteFn:                  m.podDeleteFn,
		PodContainerStartFn:          m.podContainerStartFn,
		PodContainerCreateAndStartFn: m.podContainerCreateAndStartFn,
		PodContainerRemoveFn:         m.podContainerRemoveFn,
		PodContainerRestartFn:        m.podContainerRestartFn,
	}
}

func (m *manager) UpdatePod(event *pleg.PodLifecycleEvent) {
	workCh, exists := m.workChanMap[event.ID]
	if !exists {
		return
	}
	switch event.Type {
	case pleg.ContainerNeedCreateAndStart:
		target := event.Data.(*apiObject.Container)
		workCh <- newPodContainerCreateAndStartWork(event.Pod, target)
	case pleg.ContainerNeedStart:
		workCh <- newPodContainerStartWork(event.ID, event.Data.(container.ContainerID))
	case pleg.ContainerNeedRestart:
		args := event.Data.(pleg.PodRestartContainerArgs)
		workCh <- newPodContainerRestartWork(event.Pod, args.ContainerID, args.ContainerFullName)
	}
}

func (m *manager) AddPod(pod *apiObject.Pod) {
	workCh := make(chan podWork, workChanSize)
	m.workChanMap[pod.UID()] = workCh
	worker := m.newWorker()
	workCh <- newPodCreateWork(pod)
	go worker.Run(workCh)
}

func (m *manager) DeletePod(pod *apiObject.Pod) {
	podUID := pod.UID()
	workCh := m.workChanMap[podUID]
	workCh <- newPodDeleteWork(pod)
	delete(m.workChanMap, podUID)
}

func NewPodWorkerManager(podCreateFn PodCreateFn, podDeleteFn PodDeleteFn, podContainerCreateAndStartFn PodContainerCreateAndStartFn,
	podContainerStartFn PodContainerStartFn, podContainerRemoveFn PodContainerRemoveFn, podContainerRestartFn PodContainerRestartFn) Manager {
	return &manager{
		workChanMap:                  make(workChanMap),
		podCreateFn:                  podCreateFn,
		podDeleteFn:                  podDeleteFn,
		podContainerCreateAndStartFn: podContainerCreateAndStartFn,
		podContainerStartFn:          podContainerStartFn,
		podContainerRemoveFn:         podContainerRemoveFn,
		podContainerRestartFn:        podContainerRestartFn,
	}
}
