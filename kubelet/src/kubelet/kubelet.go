package kubelet

import (
	"fmt"
	"minik8s/apiObject"
	"minik8s/kubelet/src/runtime/pleg"
	"minik8s/kubelet/src/runtime/podworker"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/kubelet/src/status"
)

type Kubelet struct {
	statusManager    status.Manager
	runtimeManager   runtime.Manager
	plegManager      pleg.Manager
	podWorkerManager podworker.Manager
}

func NewKubelet() *Kubelet {
	kl := &Kubelet{
		statusManager:  status.NewStatusManager(),
		runtimeManager: runtime.NewPodManager(),
	}
	kl.plegManager = pleg.NewPlegManager(kl.statusManager, kl.runtimeManager)
	kl.podWorkerManager = podworker.NewPodWorkerManager(kl.runtimeManager.CreatePod,
		kl.runtimeManager.PodCreateAndStartContainer,
		kl.runtimeManager.PodStartContainer,
		kl.runtimeManager.PodRemoveContainer,
		kl.runtimeManager.PodRestartContainer)
	return kl
}

func (kl *Kubelet) Run(updates <-chan *apiObject.Pod) {
	go kl.plegManager.Start()
	kl.syncLoop(updates)
}

func (kl *Kubelet) syncLoop(updates <-chan *apiObject.Pod) {
	for kl.syncLoopIteration(updates) {

	}
}

func (kl *Kubelet) syncLoopIteration(updates <-chan *apiObject.Pod) bool {
	select {
	case podUpdate := <-updates:
		fmt.Printf("podUpdate: %v\n", podUpdate)
		podUID := podUpdate.UID()

		// If pod is newly created
		if kl.statusManager.GetPod(podUID) == nil {
			kl.statusManager.UpdatePod(podUpdate.UID(), podUpdate)
			kl.podWorkerManager.AddPod(podUpdate)
		}

	case event := <-kl.plegManager.Updates():
		fmt.Printf("Receive ple: %v, data = %v\n", event, event.Data)
		kl.podWorkerManager.UpdatePod(event)
	}
	return true
}
