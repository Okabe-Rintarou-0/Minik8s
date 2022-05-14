package replicaSet

import (
	"encoding/json"
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/controller/src/cache"
	"minik8s/entity"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/listwatch"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"minik8s/util/uidutil"
	"time"
)

var logWorker = logger.Log("ReplicaSet Worker")

const timeoutSeconds = 30
const workChanSize = 5

type Worker interface {
	Run()
	SyncChannel() chan<- struct{}
	SetTarget(rs *apiObject.ReplicaSet)
}

type worker struct {
	syncCh       chan struct{}
	cacheManager cache.Manager
	target       *apiObject.ReplicaSet
}

func (w *worker) SetTarget(rs *apiObject.ReplicaSet) {
	if rs != nil {
		w.target = rs
	}
}

func (w *worker) SyncChannel() chan<- struct{} {
	return w.syncCh
}

/// TODO replace it with api-server, urgent!
// testMap is just for *TEST*, do not use it.
// We have to save pod into map, because we don't have an api-server now
// All we can do is to *Mock*
var testMap = map[types.UID]*apiObject.Pod{}

func (w *worker) addPodToApiServerForTest() {
	podTemplate := w.target.Template()
	pod := podTemplate.ToPod()
	pod.Metadata.Name = w.target.Name()
	pod.Metadata.Namespace = w.target.Namespace()
	topic := topicutil.SchedulerPodUpdateTopic()
	pod.Metadata.UID = uidutil.New()
	pod.AddLabel(runtime.KubernetesReplicaSetUIDLabel, w.target.UID())
	msg, _ := json.Marshal(entity.PodUpdate{
		Action: entity.CreateAction,
		Target: *pod,
	})
	testMap[pod.UID()] = pod
	listwatch.Publish(topic, msg)
}

func (w *worker) addPod() {
	// Add two, so we can test the case that the number of existent pods is more than replicas
	// just for test now
	w.addPodToApiServerForTest()

	// for test:
	//w.addPodToApiServerForTest()
	//w.addPodToApiServerForTest()
	//w.addPodToApiServerForTest()
}

func (w *worker) deletePodToApiServerForTest(podUID types.UID) {
	pod2Delete := testMap[podUID]
	logWorker("Pod to delete is Pod[ID = %v]", pod2Delete.UID())
	topic := topicutil.SchedulerPodUpdateTopic()
	msg, _ := json.Marshal(entity.PodUpdate{
		Action: entity.DeleteAction,
		Target: *pod2Delete,
	})
	listwatch.Publish(topic, msg)
}

func (w *worker) deletePod(podUID types.UID) {
	// just for test now
	w.deletePodToApiServerForTest(podUID)
}

func (w *worker) numRunningPods(podStatuses []*entity.PodStatus) int {
	num := 0
	for _, podStatus := range podStatuses {
		if podStatus.Lifecycle == entity.PodRunning {
			num++
		}
	}
	return num
}

func (w *worker) syncLoopIteration() bool {
	// Step 1: Get pod statues from Cache
	podStatuses := w.cacheManager.GetReplicaSetPodStatuses(w.target.UID())
	numReplicas := w.target.Replicas()

	numPods := len(podStatuses)
	numRunningPods := w.numRunningPods(podStatuses)
	diff := numPods - numReplicas
	logWorker("Syn result: diff = %d", diff)
	cpu, mem := w.calcMetrics(podStatuses)
	if diff == 0 {
		w.ready(cpu, mem)
	} else if diff > 0 {
		podUID := podStatuses[0].ID
		w.scaling(numRunningPods, cpu, mem)
		go w.deletePod(podUID)
	} else {
		w.scaling(numRunningPods, cpu, mem)
		go w.addPod()
	}

	timeout := time.NewTimer(time.Second * timeoutSeconds)
	for {
		select {
		case _, open := <-w.syncCh:
			if !open {
				return false
			}
			return true
		case <-timeout.C:
			return true
		}
	}
}

func (w *worker) syncLoop() {
	for w.syncLoopIteration() {
	}
}

func (w *worker) Run() {
	w.syncLoop()
}

func NewWorker(target *apiObject.ReplicaSet, cacheManager cache.Manager) Worker {
	return &worker{
		syncCh:       make(chan struct{}, workChanSize),
		target:       target,
		cacheManager: cacheManager,
	}
}
