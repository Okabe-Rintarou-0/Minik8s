package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/scheduler/src/filter"
	"minik8s/scheduler/src/selector"
	"minik8s/util/httputil"
	"minik8s/util/topicutil"
)

type Scheduler interface {
	Start()
	Schedule(podUpdate *entity.PodUpdate) error
}

func New() Scheduler {
	return &scheduler{
		filter.Default(),
		selector.DefaultFactory.NewSelector(selector.Random),
	}
}

type scheduler struct {
	filter   filter.Filter
	selector selector.Selector
}

// getNodesFromApiServer get nodes from api-server
func (s *scheduler) getNodesFromApiServer() (nodes []*entity.NodeStatus) {
	_ = httputil.GetAndUnmarshal(url.Prefix+url.NodeURL, &nodes)
	return
}

func (s *scheduler) getNodes() []*entity.NodeStatus {
	return s.getNodesFromApiServer()
}

func (s *scheduler) Schedule(podUpdate *entity.PodUpdate) error {
	// Step 1: Get nodes from api-server
	nodes := s.getNodes()

	if len(nodes) == 0 {
		return fmt.Errorf("no available node now")
	}

	// Step 2: Preliminary Filter
	filtered := s.filter.Filter(&podUpdate.Target, nodes)
	if len(filtered) == 0 {
		return fmt.Errorf("no suitable node now")
	}

	// Step 3: Select one node
	node := s.selector.Select(filtered)
	if node == nil {
		return fmt.Errorf("no suitable node now")
	}

	// Step 4: Prepare for the message
	nodeName := node.Hostname
	topic := topicutil.PodUpdateTopic(nodeName)
	updateMsg, err := json.Marshal(podUpdate)
	if err != nil {
		return err
	}

	// Step 5: Send msg to such node
	fmt.Printf("Send msg to %s: [%v]%v\n", topic, podUpdate.Action.String(), podUpdate.Target.Name())
	listwatch.Publish(topic, updateMsg)

	return nil
}

func (s *scheduler) parseAndSchedule(msg *redis.Message) {
	//fmt.Printf("scheduler received: %v\n", msg)
	podUpdate := &entity.PodUpdate{}
	err := json.Unmarshal([]byte(msg.Payload), podUpdate)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = s.Schedule(podUpdate); err != nil {
		fmt.Println(err.Error())
	}
}

func (s *scheduler) Start() {
	topic := topicutil.SchedulerPodUpdateTopic()

	listwatch.Watch(topic, s.parseAndSchedule)
}
