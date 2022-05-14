package replicaSet

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/controller/src/cache"
	"minik8s/controller/src/controller/hpa"
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
	logManager("Add replicaSet: %s_%s", rs.FullName(), rs.UID())
	hpa.AddRsForTest(rs)
	worker := NewWorker(rs, c.cacheManager)
	c.workers[rs.UID()] = worker
	go worker.Run()
}

///TODO just for test now, should replace with api-server later
func (c *controller) deleteReplicaSetPods(rs *apiObject.ReplicaSet) {
	logManager("Not delete the pods of rs: %s-%s", rs.FullName(), rs.UID())
	podStatuses := c.cacheManager.GetReplicaSetPodStatuses(rs.UID())
	for _, podStatus := range podStatuses {
		pod2Delete := testMap[podStatus.ID]
		logManager("Pod to delete is Pod[ID = %v]", pod2Delete.UID())
		topic := topicutil.SchedulerPodUpdateTopic()
		msg, _ := json.Marshal(entity.PodUpdate{
			Action: entity.DeleteAction,
			Target: *pod2Delete,
		})
		listwatch.Publish(topic, msg)
	}
}

func (c *controller) DeleteReplicaSet(rs *apiObject.ReplicaSet) {
	logManager("Delete replicaSet: %s_%s", rs.FullName(), rs.UID())
	if worker, exists := c.workers[rs.UID()]; exists {
		delete(c.workers, rs.UID())
		close(worker.SyncChannel())
		c.deleteReplicaSetPods(rs)
	}
}

func (c *controller) UpdateReplicaSet(rs *apiObject.ReplicaSet) {
	logManager("Update replicaSet: %s_%s", rs.FullName(), rs.UID())
	if worker, exists := c.workers[rs.UID()]; exists {
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
