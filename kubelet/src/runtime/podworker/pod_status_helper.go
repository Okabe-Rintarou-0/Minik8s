package podworker

import (
	"encoding/json"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util"
)

func (w *podWorker) pod2PodStatus(pod *apiObject.Pod) *entity.PodStatus {
	return &entity.PodStatus{
		ID:        pod.UID(),
		Name:      pod.Name(),
		Labels:    pod.Labels(),
		Namespace: pod.Namespace(),
	}
}

func (w *podWorker) runningPodStatus(pod *apiObject.Pod) *entity.PodStatus {
	podStatus := w.pod2PodStatus(pod)
	podStatus.Status = entity.Running
	return podStatus
}

func (w *podWorker) containerCreatingPodStatus(pod *apiObject.Pod) *entity.PodStatus {
	podStatus := w.pod2PodStatus(pod)
	podStatus.Status = entity.ContainerCreating
	return podStatus
}

func (w *podWorker) deletedPodStatus(pod *apiObject.Pod) *entity.PodStatus {
	podStatus := w.pod2PodStatus(pod)
	podStatus.Status = entity.Deleted
	return podStatus
}

func (w *podWorker) errorPodStatus(pod *apiObject.Pod) *entity.PodStatus {
	podStatus := w.pod2PodStatus(pod)
	podStatus.Status = entity.Error
	return podStatus
}

func (w *podWorker) publishPodStatus(podStatus *entity.PodStatus) {
	topic := util.PodStatusTopic()
	msg, _ := json.Marshal(*podStatus)
	listwatch.Publish(topic, msg)
}
