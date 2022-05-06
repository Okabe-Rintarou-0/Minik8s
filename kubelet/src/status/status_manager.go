package status

import (
	"fmt"
	"minik8s/apiObject"
	"minik8s/kubelet/src/runtime/cache"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/kubelet/src/types"
	"sync"
	"time"
)

const syncIntervalSeconds = 10

type Manager interface {
	GetPod(podUID types.UID) *apiObject.Pod
	UpdatePod(podUID types.UID, newPod *apiObject.Pod)
	DeletePod(podUID types.UID)
	GetPodStatuses() runtime.PodStatuses
	Start()
}

type manager struct {
	podCache        cache.Cache
	podStatuses     runtime.PodStatuses
	podStatusesLock sync.Mutex
	runtimeManager  runtime.Manager
}

func (m *manager) GetPod(podUID types.UID) *apiObject.Pod {
	if podVal := m.podCache.Get(podUID); podVal != nil {
		return podVal.(*apiObject.Pod)
	}
	return nil
}

func (m *manager) UpdatePod(podUID types.UID, newPod *apiObject.Pod) {
	m.podCache.Update(podUID, newPod)
}

func (m *manager) DeletePod(podUID types.UID) {
	m.podCache.Delete(podUID)
}

func (m *manager) syncWithApiServer() error {
	return nil
}

func (m *manager) syncLoop() {
	var err error
	m.podStatusesLock.Lock()
	m.podStatuses, err = m.runtimeManager.GetPodStatuses()
	m.podStatusesLock.Unlock()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = m.syncWithApiServer(); err != nil {
		fmt.Println(err.Error())
	}
}

func (m *manager) run() {
	ticker := time.Tick(time.Second * syncIntervalSeconds)
	for {
		select {
		case <-ticker:
			m.syncLoop()
		}
	}
}

func (m *manager) Start() {
	go m.run()
}

func (m *manager) GetPodStatuses() runtime.PodStatuses {
	m.podStatusesLock.Lock()
	defer m.podStatusesLock.Unlock()
	return m.podStatuses
}

func NewStatusManager(runtimeManager runtime.Manager) Manager {
	return &manager{
		podCache:       cache.Default(),
		runtimeManager: runtimeManager,
	}
}
