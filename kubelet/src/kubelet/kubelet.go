package kubelet

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/entity"
	"minik8s/kubelet/src/pleg"
	"minik8s/kubelet/src/runtime/podworker"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/kubelet/src/status"
	"minik8s/listwatch"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"os"
)

var log = logger.Log("Kubelet")

type Kubelet struct {
	statusManager    status.Manager
	runtimeManager   runtime.Manager
	plegManager      pleg.Manager
	podWorkerManager podworker.Manager

	updates chan *entity.PodUpdate
}

func New() *Kubelet {
	kl := &Kubelet{
		runtimeManager: runtime.NewPodManager(),
		updates:        make(chan *entity.PodUpdate, 20),
	}
	kl.podWorkerManager = podworker.NewPodWorkerManager(
		kl.runtimeManager.CreatePod,
		kl.runtimeManager.DeletePod,
		kl.runtimeManager.PodCreateAndStartContainer,
		kl.runtimeManager.PodStartContainer,
		kl.runtimeManager.PodRemoveContainer,
		kl.runtimeManager.PodRestartContainer,
	)
	kl.statusManager = status.NewStatusManager(
		kl.runtimeManager,
		kl.podWorkerManager.AddPod,
		kl.podWorkerManager.DeletePod,
	)
	kl.plegManager = pleg.NewPlegManager(kl.statusManager)
	return kl
}

func (kl *Kubelet) parsePodUpdate(msg *redis.Message) {
	podUpdate := &entity.PodUpdate{}
	err := json.Unmarshal([]byte(msg.Payload), podUpdate)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	log("Received pod update action: %s for %s", podUpdate.Action.String(), podUpdate.Target.UID())
	kl.updates <- podUpdate
}

func (kl *Kubelet) podUpdateTopic() (string, error) {
	hostname, err := os.Hostname()
	return topicutil.PodUpdateTopic(hostname), err
}

func (kl *Kubelet) Run() {
	topic, err := kl.podUpdateTopic()
	if err != nil {
		panic(err)
	}

	go listwatch.Watch(topic, kl.parsePodUpdate)

	kl.statusManager.Start()
	kl.plegManager.Start()
	kl.syncLoop(kl.updates)
}

func (kl *Kubelet) syncLoop(updates <-chan *entity.PodUpdate) {
	for kl.syncLoopIteration(updates) {
	}
}

func (kl *Kubelet) syncLoopIteration(updates <-chan *entity.PodUpdate) bool {
	log("Sync loop Iteration")
	select {
	case podUpdate := <-updates:
		log("Received podUpdate %v: Pod[Name = %v, UID = %v]", podUpdate.Action.String(), podUpdate.Target.FullName(), podUpdate.Target.UID())
		pod := &podUpdate.Target
		podUID := pod.UID()
		switch podUpdate.Action {
		case entity.CreateAction:
			kl.statusManager.AddPod(podUID, pod)
			kl.podWorkerManager.AddPod(pod)
		case entity.UpdateAction:
			kl.statusManager.UpdatePod(podUID, pod)
			kl.podWorkerManager.UpdatePod(pod)
		case entity.DeleteAction:
			kl.statusManager.DeletePod(podUID)
			kl.podWorkerManager.DeletePod(pod)
		}

	case event := <-kl.plegManager.Updates():
		log("Receive ple: %v, data = %v", event, event.Data)
		kl.podWorkerManager.SyncPod(event)
	}
	return true
}
