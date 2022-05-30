package replicaSet

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/controller/src/cache"
	"minik8s/entity"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/listwatch"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
)

var logManager = logger.Log("ReplicaSet Manager")
var syncSignal = struct{}{}

type Controller interface {
	Run()
	Sync(podStatus *entity.PodStatus)
	AddReplicaSet(rs *apiObject.ReplicaSet)
	DeleteReplicaSet(rs *apiObject.ReplicaSet)
	UpdateReplicaSet(rs *apiObject.ReplicaSet)
}

type controller struct {
	cacheManager cache.Manager
	workers      map[types.UID]Worker
}

func (c *controller) Sync(podStatus *entity.PodStatus) {
	needSync := podStatus.Lifecycle == entity.PodDeleted || podStatus.Lifecycle == entity.PodCreated
	if UID, exists := podStatus.Labels[runtime.KubernetesReplicaSetUIDLabel]; exists && needSync {
		if worker, stillWorking := c.workers[UID]; stillWorking {
			worker.SyncChannel() <- syncSignal
			logManager("Sync called to %s", UID)
		}
	}
}

func (c *controller) AddReplicaSet(rs *apiObject.ReplicaSet) {
	UID := rs.UID()
	logManager("Add replicaSet: %s_%s", rs.FullName(), UID)
	if _, stillWorking := c.workers[UID]; !stillWorking {
		worker := NewWorker(rs, c.cacheManager)
		c.workers[UID] = worker
		go worker.Run()
	}
}

func (c *controller) deleteReplicaSetPods(rs *apiObject.ReplicaSet) {
	logManager("Now delete the pods of rs: %s-%s", rs.FullName(), rs.UID())
	podStatuses := c.cacheManager.GetReplicaSetPodStatuses(rs.UID())
	for _, podStatus := range podStatuses {
		deletePodToApiServer(podStatus.Namespace, podStatus.Name)
	}
}

func (c *controller) DeleteReplicaSet(rs *apiObject.ReplicaSet) {
	UID := rs.UID()
	logManager("Delete replicaSet: %s_%s", rs.FullName(), UID)
	if worker, stillWorking := c.workers[UID]; stillWorking {
		close(worker.SyncChannel())
		c.deleteReplicaSetPods(rs)
		worker.Done()
		delete(c.workers, UID)
	}
}

func (c *controller) UpdateReplicaSet(rs *apiObject.ReplicaSet) {
	UID := rs.UID()
	logManager("Update replicaSet: %s_%s", rs.FullName(), UID)
	if worker, stillWorking := c.workers[UID]; stillWorking {
		worker.SetTarget(rs)
		// Sync immediately after update the rs.
		worker.SyncChannel() <- syncSignal
	}
}

func (c *controller) parseReplicaSetUpdate(msg *redis.Message) {
	replicaSetUpdate := &entity.ReplicaSetUpdate{}
	err := json.Unmarshal([]byte(msg.Payload), replicaSetUpdate)
	if err != nil {
		logManager(err.Error())
		return
	}
	rs := &replicaSetUpdate.Target

	switch replicaSetUpdate.Action {
	case entity.CreateAction:
		c.AddReplicaSet(rs)
	case entity.UpdateAction:
		c.UpdateReplicaSet(rs)
	case entity.DeleteAction:
		c.DeleteReplicaSet(rs)
	}
}

func (c *controller) Run() {
	topic := topicutil.ReplicaSetUpdateTopic()
	listwatch.Watch(topic, c.parseReplicaSetUpdate)
}

func NewController(cacheManager cache.Manager) Controller {
	return &controller{
		cacheManager: cacheManager,
		workers:      make(map[types.UID]Worker),
	}
}
