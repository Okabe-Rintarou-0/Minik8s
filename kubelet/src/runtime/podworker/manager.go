package podworker

import (
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/kubelet/src/pleg"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/util/logger"
)

var log = logger.Log("PodWorker Manager")

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
	if worker, stillWorking := m.workers[podUID]; stillWorking {
		worker.AddWork(newPodDeleteWork(newPod))
		worker.AddWork(newPodCreateWork(newPod))
	}
}

func (m *manager) SyncPod(event *pleg.PodLifecycleEvent) {
	worker, exists := m.workers[event.ID]
	if !exists {
		return
	}
	switch event.Type {
	case pleg.ContainerNeedCreateAndStart:
		target := event.Data.(*apiObject.Container)
		worker.AddWork(newPodContainerCreateAndStartWork(event.Pod, target))
	case pleg.ContainerNeedStart:
		worker.AddWork(newPodContainerStartWork(event.ID, event.Data.(container.ID)))
	case pleg.ContainerNeedRestart:
		args := event.Data.(pleg.PodRestartContainerArgs)
		worker.AddWork(newPodContainerRestartWork(event.Pod, args.ContainerID, args.ContainerFullName))
	}
}

func (m *manager) AddPod(pod *apiObject.Pod) {
	worker := m.newWorker()
	m.workers[pod.UID()] = worker
	worker.AddWork(newPodCreateWork(pod))
	go worker.Run()
}

func (m *manager) DeletePod(pod *apiObject.Pod) {
	log("Delete pod[ID = %s]", pod.UID())
	podUID := pod.UID()
	if worker, stillWorking := m.workers[podUID]; stillWorking {
		worker.AddWork(newPodDeleteWork(pod))
		worker.Done()
		delete(m.workers, podUID)
	}
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
