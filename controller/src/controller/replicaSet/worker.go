package replicaSet

import (
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/controller/src/cache"
	"minik8s/entity"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/util/httputil"
	"minik8s/util/logger"
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

func (w *worker) addPodToApiServer() {
	podTemplate := w.target.Template()
	pod := podTemplate.ToPod()
	pod.Metadata.Name = "ReplicaSet-" + w.target.Name() + "-" + uidutil.New()
	pod.Metadata.Namespace = w.target.Namespace()
	pod.AddLabel(runtime.KubernetesReplicaSetUIDLabel, w.target.UID())

	URL := url.Prefix + url.PodURL
	if _, err := httputil.PostJson(URL, pod); err != nil {
		logger.Error(err.Error())
	}
}

func (w *worker) addPod() {
	w.addPodToApiServer()
}

func (w *worker) deletePod(namespace, name string) {
	deletePodToApiServer(namespace, name)
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
		podToDelete := podStatuses[0]
		w.scaling(numRunningPods, cpu, mem)
		go w.deletePod(podToDelete.Namespace, podToDelete.Name)
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
