package status

import (
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/entity"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/listwatch"
	"minik8s/util/cache"
	"minik8s/util/logger"
	"minik8s/util/parseutil"
	"minik8s/util/topicutil"
	"minik8s/util/utility"
	"minik8s/util/wait"
	"os"
	"sync"
	"time"
)

var log = logger.Log("Status Manager")

const syncPeriod = 10 * time.Second

type Manager interface {
	GetPod(podUID types.UID) *apiObject.Pod
	AddPod(podUID types.UID, pod *apiObject.Pod)
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

func (m *manager) AddPod(podUID types.UID, pod *apiObject.Pod) {
	m.podCache.Add(podUID, pod)
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

//TODO call REST api of api-server? or just publish it to the topic and api-server also watches it?
func (m *manager) publishPodStatus(podStatuses runtime.PodStatuses) {
	topic := topicutil.PodStatusTopic()
	for _, podStatus := range podStatuses {
		e := podStatus.ToEntity()
		m.addLabels(e)
		log("Publish Pod[ID = %v, cpu = %v, mem = %v]", e.ID, e.CpuPercent, e.MemPercent)
		listwatch.Publish(topic, parseutil.MarshalAny(e))
	}
}

func (m *manager) publishNodeStatus(podStatuses runtime.PodStatuses) {
	hostname, _ := os.Hostname()
	cpu, mem := utility.GetCpuAndMemoryUsage()
	nodeStatus := &entity.NodeStatus{
		Hostname:   hostname,
		Lifecycle:  entity.NodeReady,
		Error:      "",
		CpuPercent: cpu,
		MemPercent: mem,
		NumPods:    len(podStatuses),
		SyncTime:   time.Now(),
	}
	topic := topicutil.NodeStatusTopic()
	listwatch.Publish(topic, parseutil.MarshalAny(nodeStatus))
	log("Publish Node Status[Host: %v, cpu: %v, mem: %v, pods: %v]", hostname, cpu, mem, nodeStatus.NumPods)
}

func (m *manager) syncWithApiServer(podStatuses runtime.PodStatuses) error {
	m.publishPodStatus(podStatuses)
	m.publishNodeStatus(podStatuses)
	return nil
}

func (m *manager) addLabels(podStatus *entity.PodStatus) {
	v := m.podCache.Get(podStatus.ID)
	if v == nil {
		return
	}
	podStatus.Labels = v.(*apiObject.Pod).Labels().DeepCopy()
}

func (m *manager) syncLoopIteration() {
	podStatuses, err := m.runtimeManager.GetPodStatuses()
	if err != nil {
		log(err.Error())
		return
	}
	m.podStatusesLock.Lock()
	m.podStatuses = podStatuses
	m.podStatusesLock.Unlock()

	if err = m.syncWithApiServer(podStatuses); err != nil {
		log(err.Error())
	}
}

func (m *manager) syncLoop() {
	wait.Period(syncPeriod, syncPeriod, m.syncLoopIteration)
}

func (m *manager) Start() {
	go m.syncLoop()
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
