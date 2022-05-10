package podworker

import (
	"fmt"
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/kubelet/src/pleg"
	"minik8s/kubelet/src/runtime/container"
)

type workers map[types.UID]*podWorker

type Manager interface {
	AddPod(pod *apiObject.Pod)
	DeletePod(pod *apiObject.Pod)
	UpdatePod(newPod *apiObject.Pod)
	SyncPod(event *pleg.PodLifecycleEvent)
}

type manager struct {
	workers                      workers
	podCreateFn                  PodCreateFn
	podDeleteFn                  PodDeleteFn
	podContainerCreateAndStartFn PodContainerCreateAndStartFn
	podContainerStartFn          PodContainerStartFn
	podContainerRemoveFn         PodContainerRemoveFn
	podContainerRestartFn        PodContainerRestartFn
}

func (m *manager) newWorker() *podWorker {
	return newWorker(m.podCreateFn,
		m.podDeleteFn,
		m.podContainerCreateAndStartFn,
		m.podContainerStartFn,
		m.podContainerRemoveFn,
		m.podContainerRestartFn,
	)
}

func (m *manager) UpdatePod(newPod *apiObject.Pod) {
	podUID := newPod.UID()
	workCh := m.workers[podUID].WorkChannel()
	workCh <- newPodDeleteWork(newPod)
	workCh <- newPodCreateWork(newPod)
}

func (m *manager) SyncPod(event *pleg.PodLifecycleEvent) {
	worker, exists := m.workers[event.ID]
	if !exists {
		return
	}
	workCh := worker.WorkChannel()
	switch event.Type {
	case pleg.ContainerNeedCreateAndStart:
		target := event.Data.(*apiObject.Container)
		workCh <- newPodContainerCreateAndStartWork(event.Pod, target)
	case pleg.ContainerNeedStart:
		workCh <- newPodContainerStartWork(event.ID, event.Data.(container.ID))
	case pleg.ContainerNeedRestart:
		args := event.Data.(pleg.PodRestartContainerArgs)
		workCh <- newPodContainerRestartWork(event.Pod, args.ContainerID, args.ContainerFullName)
	}
}

func (m *manager) AddPod(pod *apiObject.Pod) {
	worker := m.newWorker()
	m.workers[pod.UID()] = worker
	workCh := worker.WorkChannel()
	workCh <- newPodCreateWork(pod)
	go worker.Run()
}

func (m *manager) DeletePod(pod *apiObject.Pod) {
	fmt.Printf("[PodWorkerManager]: delete pod %s\n", pod.UID())
	podUID := pod.UID()
	workCh := m.workers[podUID].WorkChannel()
	workCh <- newPodDeleteWork(pod)
	delete(m.workers, podUID)
	close(workCh)
}

func NewPodWorkerManager(podCreateFn PodCreateFn, podDeleteFn PodDeleteFn, podContainerCreateAndStartFn PodContainerCreateAndStartFn,
	podContainerStartFn PodContainerStartFn, podContainerRemoveFn PodContainerRemoveFn, podContainerRestartFn PodContainerRestartFn) Manager {
	return &manager{
		workers:                      make(workers),
		podCreateFn:                  podCreateFn,
		podDeleteFn:                  podDeleteFn,
		podContainerCreateAndStartFn: podContainerCreateAndStartFn,
		podContainerStartFn:          podContainerStartFn,
		podContainerRemoveFn:         podContainerRemoveFn,
		podContainerRestartFn:        podContainerRestartFn,
	}
}
