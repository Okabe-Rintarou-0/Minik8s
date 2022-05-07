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

const workChanSize = 5

type controller struct {
	cacheManager cache.Manager
	syncChannels map[types.UID]chan struct{}
}

func (c *controller) Sync(podStatus *entity.PodStatus) {
	needSync := podStatus.Status == entity.Deleted || podStatus.Status == entity.Running
	if UID, exists := podStatus.Labels[runtime.KubernetesReplicaSetUIDLabel]; exists && needSync {
		c.syncChannels[UID] <- struct{}{}
		fmt.Printf("Sync called to %s\n", UID)
	}
}

func (c *controller) AddReplicaSet(rs *apiObject.ReplicaSet) {
	fmt.Println("add rs!")
	syncCh := make(chan struct{}, workChanSize)
	c.syncChannels[rs.UID()] = syncCh
	worker := NewWorker(syncCh, rs, c.cacheManager)
	go worker.Run()
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

	case entity.DeleteAction:

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
		syncChannels: make(map[types.UID]chan struct{}),
	}
}
