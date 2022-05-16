package scheduler

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/scheduler/src/selector"
	"testing"
)

func readPod(podPath string) *apiObject.Pod {
	content, _ := ioutil.ReadFile(podPath)
	pod := apiObject.Pod{}
	_ = yaml.Unmarshal(content, &pod)
	return &pod
}

func TestGetNodes(t *testing.T) {
	sch := scheduler{selector: selector.DefaultFactory.NewSelector(selector.Random)}
	fmt.Println(sch.getNodes()[0])
}

func TestScheduler(t *testing.T) {
	s := New()
	err := s.Schedule(&entity.PodUpdate{
		Action: entity.CreateAction,
		Target: *testPod,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
