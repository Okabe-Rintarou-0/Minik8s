package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/scheduler/src/selector"
	"minik8s/util"
	"os"
)

type Scheduler interface {
	Start()
	Schedule(podUpdate *entity.PodUpdate) error
}

func New() Scheduler {
	return &scheduler{selector.New()}
}

type scheduler struct {
	selector selector.Selector
}

func getTestNodes() []*apiObject.Node {
	var nodes []*apiObject.Node
	node := &apiObject.Node{}
	node.ApiVersion = "v1"
	node.Kind = "Node"
	hostname, _ := os.Hostname()
	node.Metadata.Name = hostname
	node.Metadata.Namespace = "default"
	nodes = append(nodes, node)
	return nodes
}

// getNodes should get nodes from api server
func (s *scheduler) getNodes() []*apiObject.Node {
	// for test now
	return getTestNodes()
}

func (s *scheduler) Schedule(podUpdate *entity.PodUpdate) error {
	//fmt.Printf("Schedule %v\n", podUpdate)
	// Step 1: Get nodes from api-server
	nodes := s.getNodes()

	if len(nodes) == 0 {
		return fmt.Errorf("no available node now")
	}

	// Step 2: Select one node
	node := s.selector.Select(nodes)
	if node == nil {
		return fmt.Errorf("no suitable node now")
	}

	// Step 3: Prepare for the message
	nodeName := node.Metadata.Name
	topic := util.PodUpdateTopic(nodeName)
	updateMsg, err := json.Marshal(*podUpdate)
	if err != nil {
		return err
	}

	// Step 4: Send msg to such node
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
	topic := util.SchedulerPodUpdateTopic()

	listwatch.Watch(topic, s.parseAndSchedule)
}
