package kubelet

import (
	"fmt"
	"minik8s/apiObject"
	"minik8s/kubelet/src/runtime/pleg"
	"minik8s/kubelet/src/runtime/pod"
	"minik8s/kubelet/src/status"
)

type Kubelet struct {
	statusManager status.Manager
	podManager    pod.Manager
	plegManager   pleg.Manager
}

func NewKubelet() *Kubelet {
	kl := &Kubelet{
		statusManager: status.NewStatusManager(),
		podManager:    pod.NewPodManager(),
	}
	kl.plegManager = pleg.NewPlegManager(kl.statusManager, kl.podManager)
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
		kl.statusManager.UpdatePod(podUpdate.UID(), podUpdate)
		err := kl.podManager.CreatePod(podUpdate)
		if err != nil {
			return false
		}
	case lifecycleEvent := <-kl.plegManager.Updates():
		fmt.Printf("Receive ple: %v, data = %v\n", lifecycleEvent, lifecycleEvent.Data)
	default:
	}
	return true
}
