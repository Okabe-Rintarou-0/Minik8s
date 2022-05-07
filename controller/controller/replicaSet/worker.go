package replicaSet

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/controller/controller/cache"
	"minik8s/entity"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/listwatch"
	"minik8s/util"
	"time"
)

const timeoutSeconds = 30

type Worker interface {
	Run()
}

type worker struct {
	syncCh       chan struct{}
	cacheManager cache.Manager
	target       *apiObject.ReplicaSet
}

// testMap is just for *TEST*, do not use it.
// We have to save pod into map, because we don't have an api-server now
// All we can do is to *Mock*
var testMap = map[types.UID]*apiObject.Pod{}

func (w *worker) addPodToApiServerForTest() {
	topic := util.SchedulerPodUpdateTopic()
	content, _ := ioutil.ReadFile("./testPod.yaml")
	pod := apiObject.Pod{}
	pod.Metadata.UID = uuid.NewV4().String()
	pod.AddLabel(runtime.KubernetesReplicaSetUIDLabel, w.target.UID())
	_ = yaml.Unmarshal(content, &pod)
	fmt.Println("Sending test pod", pod)
	msg, _ := json.Marshal(entity.PodUpdate{
		Action: entity.CreateAction,
		Target: pod,
	})
	testMap[pod.UID()] = &pod
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
	fmt.Printf("Pod 2 delete is %v\n", pod2Delete)
	topic := util.SchedulerPodUpdateTopic()
	content, _ := ioutil.ReadFile("./testPod.yaml")
	_ = yaml.Unmarshal(content, pod2Delete)
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

func (w *worker) syncLoopIteration() bool {
	// Step 1: Get pod statues from Cache
	podStatuses := w.cacheManager.GetReplicaSetPodStatuses(w.target.UID())
	numReplicas := w.target.Replicas()

	diff := len(podStatuses) - numReplicas
	fmt.Printf("Syn result: diff = %d\n", diff)
	if diff > 0 {
		podUID := podStatuses[0].ID
		go w.deletePod(podUID)
	} else if diff < 0 {
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

// syncLoop loops every syncLoopIntervalSeconds seconds
func (w *worker) syncLoop() {
	for w.syncLoopIteration() {
	}
}

func (w *worker) Run() {
	w.syncLoop()
}

func NewWorker(syncCh chan struct{}, target *apiObject.ReplicaSet, cacheManager cache.Manager) Worker {
	return &worker{
		syncCh:       syncCh,
		target:       target,
		cacheManager: cacheManager,
	}
}
