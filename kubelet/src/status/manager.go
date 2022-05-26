package status

import (
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/listwatch"
	"minik8s/util/cache"
	"minik8s/util/httputil"
	"minik8s/util/logger"
	"minik8s/util/netutil"
	"minik8s/util/parseutil"
	"minik8s/util/topicutil"
	"minik8s/util/utility"
	"minik8s/util/wait"
	"strings"
	"sync"
	"time"
)

var log = logger.Log("Status Manager")

const (
	syncPeriod     = 10 * time.Second
	fullSyncPeriod = 30 * time.Second
)

type FullSyncAddHook func(pod *apiObject.Pod)
type FullSyncDeleteHook func(pod *apiObject.Pod)

type Manager interface {
	GetPod(podUID types.UID) *apiObject.Pod
	AddPod(podUID types.UID, pod *apiObject.Pod)
	UpdatePod(podUID types.UID, newPod *apiObject.Pod)
	DeletePod(podUID types.UID)
	GetPodStatuses() runtime.PodStatuses
	Start()
}

type manager struct {
	podCache           cache.Cache
	podStatuses        runtime.PodStatuses
	podStatusesLock    sync.Mutex
	runtimeManager     runtime.Manager
	fullSyncAddHook    FullSyncAddHook
	fullSyncDeleteHook FullSyncDeleteHook
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
		if m.addInfo(e) {
			log("Publish Pod %s/%s[ID = %v, cpu = %v, mem = %v]", e.Namespace, e.Name, e.ID, e.CpuPercent, e.MemPercent)
			listwatch.Publish(topic, parseutil.MarshalAny(e))
		}
	}
}

func (m *manager) publishNodeStatus(podStatuses runtime.PodStatuses) {
	hostname := netutil.Hostname()
	cpu, mem := utility.GetCpuAndMemoryUsage()
	nodeStatus := &entity.NodeStatus{
		Hostname:   hostname,
		Namespace:  "default",
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

func (m *manager) addInfo(podStatus *entity.PodStatus) bool {
	v := m.podCache.Get(podStatus.ID)
	if v == nil {
		return false
	}
	pod := v.(*apiObject.Pod)
	podStatus.Labels = pod.Labels().DeepCopy()
	podStatus.Namespace = pod.Namespace()
	podStatus.Name = pod.Name()
	podStatus.Ip = pod.ClusterIp()
	return true
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

func (m *manager) computePodAction(cachedPodMap map[string]interface{}, pods []*apiObject.Pod) (toAdd, toDelete []*apiObject.Pod) {
	for _, pod := range pods {
		if _, exists := cachedPodMap[pod.UID()]; !exists {
			toAdd = append(toAdd, pod)
		} else {
			delete(cachedPodMap, pod.UID())
		}
	}

	for _, cached := range cachedPodMap {
		toDelete = append(toDelete, cached.(*apiObject.Pod))
	}

	return
}

func (m *manager) handleAdd(toAdd []*apiObject.Pod) {
	for _, podToAdd := range toAdd {
		m.podCache.Add(podToAdd.UID(), podToAdd)
		m.fullSyncAddHook(podToAdd)
	}
}

func (m *manager) handleDelete(toDelete []*apiObject.Pod) {
	for _, podToDelete := range toDelete {
		m.podCache.Delete(podToDelete.UID())
		m.fullSyncDeleteHook(podToDelete)
	}
}

func (m *manager) fullSyncLoopIteration() {
	log("Full sync pod apiObject")
	URL := url.Prefix + strings.Replace(url.PodURLWithSpecifiedNode, ":node", netutil.Hostname(), 1)
	var pods []*apiObject.Pod
	if err := httputil.GetAndUnmarshal(URL, &pods); err == nil {
		podMap := m.podCache.ToMap()
		toAdd, toDelete := m.computePodAction(podMap, pods)

		m.handleAdd(toAdd)
		m.handleDelete(toDelete)
	} else {
		logger.Error(err.Error())
	}
}

func (m *manager) syncLoop() {
	go wait.Period(syncPeriod, syncPeriod, m.syncLoopIteration)
	go wait.Period(fullSyncPeriod, fullSyncPeriod, m.fullSyncLoopIteration)
}

func (m *manager) Start() {
	m.syncLoop()
}

func (m *manager) GetPodStatuses() runtime.PodStatuses {
	m.podStatusesLock.Lock()
	defer m.podStatusesLock.Unlock()
	return m.podStatuses
}

func NewStatusManager(runtimeManager runtime.Manager, fullSyncAddHook FullSyncAddHook, fullSyncDeleteHook FullSyncDeleteHook) Manager {
	return &manager{
		podCache:           cache.Default(),
		runtimeManager:     runtimeManager,
		fullSyncAddHook:    fullSyncAddHook,
		fullSyncDeleteHook: fullSyncDeleteHook,
	}
}
