package kubelet

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"minik8s/apiObject"
	"minik8s/kubelet/src/listwatch"
	"minik8s/kubelet/src/pleg"
	"minik8s/kubelet/src/runtime/podworker"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/kubelet/src/status"
)

const (
	CreateAction PodUpdateAction = iota
	DeleteAction
	UpdateAction
)

const PodUpdateTopic = "PodUpdate"

type PodUpdateAction byte

type PodUpdate struct {
	Action PodUpdateAction
	Target *apiObject.Pod
}

type Kubelet struct {
	statusManager    status.Manager
	runtimeManager   runtime.Manager
	plegManager      pleg.Manager
	podWorkerManager podworker.Manager

	updates chan *PodUpdate
}

func NewKubelet() *Kubelet {
	kl := &Kubelet{
		statusManager:  status.NewStatusManager(),
		runtimeManager: runtime.NewPodManager(),

		updates: make(chan *PodUpdate),
	}
	kl.plegManager = pleg.NewPlegManager(kl.statusManager, kl.runtimeManager)
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
	podUpdate := &PodUpdate{}
	err := json.Unmarshal([]byte(msg.Payload), podUpdate)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	kl.updates <- podUpdate
}

func (kl *Kubelet) Run() {
	go listwatch.Watch(PodUpdateTopic, kl.parsePodUpdate)

	kl.plegManager.Start()
	kl.syncLoop(kl.updates)
}

func (kl *Kubelet) syncLoop(updates <-chan *PodUpdate) {
	for kl.syncLoopIteration(updates) {

	}
}

func (kl *Kubelet) syncLoopIteration(updates <-chan *PodUpdate) bool {
	select {
	case podUpdate := <-updates:
		fmt.Printf("Received podUpdate %v: %v\n", podUpdate.Action, podUpdate.Target)
		pod := podUpdate.Target
		podUID := pod.UID()
		switch podUpdate.Action {
		case CreateAction:
			// If pod is newly created
			if kl.statusManager.GetPod(podUID) == nil {
				kl.statusManager.UpdatePod(podUID, pod)
				kl.podWorkerManager.AddPod(pod)
			}
		case UpdateAction:
			kl.statusManager.UpdatePod(podUID, pod)
			kl.podWorkerManager.UpdatePod(pod)
		case DeleteAction:
			kl.statusManager.DeletePod(podUID)
			kl.podWorkerManager.DeletePod(pod)
		}

	case event := <-kl.plegManager.Updates():
		fmt.Printf("Receive ple: %v, data = %v\n", event, event.Data)
		kl.podWorkerManager.SyncPod(event)
	}
	return true
}
