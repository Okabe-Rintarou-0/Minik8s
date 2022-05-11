package kubelet

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"minik8s/entity"
	"minik8s/kubelet/src/pleg"
	"minik8s/kubelet/src/runtime/podworker"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/kubelet/src/status"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"os"
)

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
	kl.statusManager = status.NewStatusManager(kl.runtimeManager)
	kl.plegManager = pleg.NewPlegManager(kl.statusManager)
	kl.podWorkerManager = podworker.NewPodWorkerManager(
		kl.runtimeManager.CreatePod,
		kl.runtimeManager.DeletePod,
		kl.runtimeManager.PodCreateAndStartContainer,
		kl.runtimeManager.PodStartContainer,
		kl.runtimeManager.PodRemoveContainer,
		kl.runtimeManager.PodRestartContainer,
	)
	return kl
}

func (kl *Kubelet) parsePodUpdate(msg *redis.Message) {
	podUpdate := &entity.PodUpdate{}
	err := json.Unmarshal([]byte(msg.Payload), podUpdate)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Kl received pod update action: %s for %s\n", podUpdate.Action.String(), podUpdate.Target.UID())
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
	fmt.Println("Sync loop")
	select {
	case podUpdate := <-updates:
		fmt.Printf("Received podUpdate %v: %v\n", podUpdate.Action.String(), podUpdate.Target)
		pod := &podUpdate.Target
		podUID := pod.UID()
		switch podUpdate.Action {
		case entity.CreateAction:
			// If pod is newly created
			if kl.statusManager.GetPod(podUID) == nil {
				kl.statusManager.UpdatePod(podUID, pod)
				kl.podWorkerManager.AddPod(pod)
			}
		case entity.UpdateAction:
			kl.statusManager.UpdatePod(podUID, pod)
			kl.podWorkerManager.UpdatePod(pod)
		case entity.DeleteAction:
			kl.statusManager.DeletePod(podUID)
			kl.podWorkerManager.DeletePod(pod)
		}

	case event := <-kl.plegManager.Updates():
		fmt.Printf("Receive ple: %v, data = %v\n", event, event.Data)
		kl.podWorkerManager.SyncPod(event)
	}
	return true
}
