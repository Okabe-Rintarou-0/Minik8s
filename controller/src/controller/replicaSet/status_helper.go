package replicaSet

import (
	"encoding/json"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"time"
)

func (w *worker) replicaSet2ReplicaSetStatus(replicaSet *apiObject.ReplicaSet) *entity.ReplicaSetStatus {
	return &entity.ReplicaSetStatus{
		ID:          replicaSet.UID(),
		Name:        replicaSet.Name(),
		Namespace:   replicaSet.Namespace(),
		Labels:      replicaSet.Labels(),
		NumReplicas: replicaSet.Replicas(),
		SyncTime:    time.Now(),
	}
}

func (w *worker) calcMetrics(podStatuses []*entity.PodStatus) (cpu, mem float64) {
	cpu = 0.0
	mem = 0.0
	for _, podStatus := range podStatuses {
		cpu += podStatus.CpuPercent
		mem += podStatus.MemPercent
	}
	return
}

func (w *worker) ready(cpu, mem float64) {
	replicaSetStatus := w.replicaSet2ReplicaSetStatus(w.target)
	replicaSetStatus.Lifecycle = entity.ReplicaSetReady
	replicaSetStatus.NumReady = replicaSetStatus.NumReplicas
	replicaSetStatus.CpuPercent = cpu
	replicaSetStatus.MemPercent = mem
	publishReplicaSetStatus(replicaSetStatus)
}

func (w *worker) scaling(numRunningPods int, cpu, mem float64) {
	replicaSetStatus := w.replicaSet2ReplicaSetStatus(w.target)
	replicaSetStatus.Lifecycle = entity.ReplicaSetScaling
	replicaSetStatus.NumReady = numRunningPods
	replicaSetStatus.CpuPercent = cpu
	replicaSetStatus.MemPercent = mem
	publishReplicaSetStatus(replicaSetStatus)
}

func (w *worker) deleted() {
	replicaSetStatus := w.replicaSet2ReplicaSetStatus(w.target)
	replicaSetStatus.Lifecycle = entity.ReplicaSetDeleted
	replicaSetStatus.NumReady = 0
	publishReplicaSetStatus(replicaSetStatus)
}

func (w *worker) error() {
	replicaSetStatus := w.replicaSet2ReplicaSetStatus(w.target)
	replicaSetStatus.Lifecycle = entity.ReplicaSetError
	publishReplicaSetStatus(replicaSetStatus)
}

func publishReplicaSetStatus(replicaSetStatus *entity.ReplicaSetStatus) {
	topic := topicutil.ReplicaSetStatusTopic()
	msg, _ := json.Marshal(replicaSetStatus)
	//fmt.Printf("publish: %v\n", replicaSetStatus)
	listwatch.Publish(topic, msg)
}
