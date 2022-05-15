package podworker

import (
	"encoding/json"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
)

func (w *podWorker) pod2PodStatus(pod *apiObject.Pod) *entity.PodStatus {
	return &entity.PodStatus{
		ID:        pod.UID(),
		Name:      pod.Name(),
		Labels:    pod.Labels(),
		Namespace: pod.Namespace(),
	}
}

func (w *podWorker) created(pod *apiObject.Pod) {
	podStatus := w.pod2PodStatus(pod)
	podStatus.Lifecycle = entity.PodCreated
	publishPodStatus(podStatus)
}

func (w *podWorker) running(pod *apiObject.Pod) {
	podStatus := w.pod2PodStatus(pod)
	podStatus.Lifecycle = entity.PodRunning
	publishPodStatus(podStatus)
}

func (w *podWorker) containerCreating(pod *apiObject.Pod) {
	podStatus := w.pod2PodStatus(pod)
	podStatus.Lifecycle = entity.PodContainerCreating
	publishPodStatus(podStatus)
}

func (w *podWorker) deleted(pod *apiObject.Pod) {
	podStatus := w.pod2PodStatus(pod)
	podStatus.Lifecycle = entity.PodDeleted
	publishPodStatus(podStatus)
}

func (w *podWorker) error(pod *apiObject.Pod) {
	podStatus := w.pod2PodStatus(pod)
	podStatus.Lifecycle = entity.PodError
	publishPodStatus(podStatus)
}

func publishPodStatus(podStatus *entity.PodStatus) {
	topic := topicutil.PodStatusTopic()
	msg, _ := json.Marshal(podStatus)
	listwatch.Publish(topic, msg)
}
