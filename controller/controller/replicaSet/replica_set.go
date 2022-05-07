package replicaSet

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/controller/controller/cache"
	"minik8s/entity"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/listwatch"
	"minik8s/util"
)

var syncSignal = struct{}{}

type controller struct {
	cacheManager cache.Manager
	workers      map[types.UID]Worker
}

func (c *controller) Sync(podStatus *entity.PodStatus) {
	needSync := podStatus.Status == entity.Deleted || podStatus.Status == entity.Running
	if UID, exists := podStatus.Labels[runtime.KubernetesReplicaSetUIDLabel]; exists && needSync {
		if worker, stillWorking := c.workers[UID]; stillWorking {
			worker.SyncChannel() <- syncSignal
			fmt.Printf("Sync called to %s\n", UID)
		}
	}
}

func (c *controller) AddReplicaSet(rs *apiObject.ReplicaSet) {
	fmt.Printf("Add replicaSet: %s_%s\n", rs.FullName(), rs.UID())
	worker := NewWorker(rs, c.cacheManager)
	c.workers[rs.UID()] = worker
	go worker.Run()
}

///TODO just for test now, should replace with api-server later
func (c *controller) deleteReplicaSetPods(rs *apiObject.ReplicaSet) {
	fmt.Printf("Not delete the pods of rs: %s-%s\n", rs.FullName(), rs.UID())
	podStatuses := c.cacheManager.GetReplicaSetPodStatuses(rs.UID())
	for _, podStatus := range podStatuses {
		pod2Delete := testMap[podStatus.ID]
		fmt.Printf("Pod 2 delete is %v\n", pod2Delete)
		topic := util.SchedulerPodUpdateTopic()
		msg, _ := json.Marshal(entity.PodUpdate{
			Action: entity.DeleteAction,
			Target: *pod2Delete,
		})
		listwatch.Publish(topic, msg)
	}
}

func (c *controller) DeleteReplicaSet(rs *apiObject.ReplicaSet) {
	fmt.Printf("Delete replicaSet: %s_%s\n", rs.FullName(), rs.UID())
	if worker, exists := c.workers[rs.UID()]; exists {
		delete(c.workers, rs.UID())
		close(worker.SyncChannel())
		c.deleteReplicaSetPods(rs)
	}
}

func (c *controller) UpdateReplicaSet(rs *apiObject.ReplicaSet) {
	fmt.Printf("Update replicaSet: %s_%s\n", rs.FullName(), rs.UID())
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
		fmt.Println(err.Error())
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
	topic := util.ReplicaSetUpdateTopic()
	listwatch.Watch(topic, c.parseReplicaSetUpdate)
}

type Controller interface {
	Run()
	Sync(podStatus *entity.PodStatus)
}

func NewController(cacheManager cache.Manager) Controller {
	return &controller{
		cacheManager: cacheManager,
		workers:      make(map[types.UID]Worker),
	}
}
