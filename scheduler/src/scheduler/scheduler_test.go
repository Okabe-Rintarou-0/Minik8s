package scheduler

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"testing"
	"time"
)

func readPod(podPath string) *apiObject.Pod {
	content, _ := ioutil.ReadFile(podPath)
	pod := apiObject.Pod{}
	_ = yaml.Unmarshal(content, &pod)
	return &pod
}

func TestScheduler(t *testing.T) {
	pod := readPod("../../test/testPod.yaml")
	pod.Metadata.UID = uuid.NewV4().String()
	createAct := entity.PodUpdate{
		Action: entity.CreateAction,
		Target: *pod,
	}

	deleteAct := entity.PodUpdate{
		Action: entity.DeleteAction,
		Target: *pod,
	}
	createMsg, err := json.Marshal(createAct)
	if err != nil {
		fmt.Println(err.Error())
	}

	deleteMsg, err := json.Marshal(deleteAct)
	if err != nil {
		fmt.Println(err.Error())
	}

	pod2 := readPod("../../test/testPod2.yaml")
	pod2.Metadata.UID = pod.UID()
	fmt.Println(pod2)
	updateAct := entity.PodUpdate{
		Action: entity.UpdateAction,
		Target: *pod2,
	}
	updateMsg, err := json.Marshal(updateAct)
	if err != nil {
		fmt.Println(err.Error())
	}

	topic := topicutil.SchedulerPodUpdateTopic()
	// after 5s, create the pod
	// after 15s, update the pod
	// after 50s, delete the pod
	go func() {
		createTimer := time.NewTimer(time.Second * 5)
		updateTimer := time.NewTimer(time.Second * 15)
		deleteTimer := time.NewTimer(time.Second * 50)
		for i := 0; i < 3; i++ {
			select {
			case <-updateTimer.C:
				fmt.Println("Now update the pod")
				listwatch.Publish(topic, updateMsg)
			case <-deleteTimer.C:
				fmt.Println("Now delete the pod")
				listwatch.Publish(topic, deleteMsg)
			case <-createTimer.C:
				fmt.Println("Now create the pod")
				listwatch.Publish(topic, createMsg)
			}
		}
	}()

	scheduler := New()
	scheduler.Start()
}
