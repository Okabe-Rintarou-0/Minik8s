package podworker

import (
	"minik8s/apiObject"
	"minik8s/kubelet/src/types"
)

const (
	workChanSize = 5
)

type workChan = <-chan podWork

type workChanMap map[types.UID]workChan

type Manager struct {
	workChanMap workChanMap
}

func (m *Manager) AddPod(pod *apiObject.Pod) {
	workCh := make(workChan, workChanSize)
	m.workChanMap[pod.UID()] = workCh
	worker := &podWorker{}
	go worker.Run(workCh)
}
