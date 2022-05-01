package kubelet

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/kubelet/src/listwatch"
	"minik8s/kubelet/src/podutil"
	"testing"
	"time"
)

func readPod(podPath string) *apiObject.Pod {
	content, _ := ioutil.ReadFile(podPath)
	pod := apiObject.Pod{}
	_ = yaml.Unmarshal(content, &pod)
	return &pod
}

func TestKubelet(t *testing.T) {
	pod := readPod("../../test/testPod.yaml")
	pod.Metadata.UID = uuid.NewV4().String()
	createAct := PodUpdate{
		Action: CreateAction,
		Target: pod,
	}
	deleteAct := PodUpdate{
		Action: DeleteAction,
		Target: pod,
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
	updateAct := PodUpdate{
		Action: UpdateAction,
		Target: pod2,
	}
	updateMsg, err := json.Marshal(updateAct)
	if err != nil {
		fmt.Println(err.Error())
	}

	// after 5s, create the pod
	// after 15s, update the pod
	// after 50s, delete the pod
	go func() {
		createTimer := time.NewTimer(time.Second * 5)
		updateTimer := time.NewTimer(time.Second * 50)
		deleteTimer := time.NewTimer(time.Second * 100)
		for i := 0; i < 3; i++ {
			select {
			case <-deleteTimer.C:
				fmt.Println("Now delete the pod")
				listwatch.Publish(PodUpdateTopic, deleteMsg)
			case <-updateTimer.C:
				fmt.Println("Now update the pod")
				listwatch.Publish(PodUpdateTopic, updateMsg)
			case <-createTimer.C:
				fmt.Println("Now create the pod")
				listwatch.Publish(PodUpdateTopic, createMsg)
			}
		}
	}()

	kl := NewKubelet()
	kl.Run()
}

func TestParse(t *testing.T) {
	str := "/k8s_POD_test_test_87653ecc-04c3-4b32-b4e3-94346b968ede_0"
	succ, containerName, podName, namespace, uid, count := podutil.ParseContainerFullName(str)
	fmt.Println(succ, containerName, podName, namespace, uid, count)
}

func TestCreatePodWithoutSpecifiedPort(t *testing.T) {
	pod := readPod("../../test/testPodWithoutSpecifiedPort.yaml")
	pod.Metadata.UID = uuid.NewV4().String()
	createAct := PodUpdate{
		Action: CreateAction,
		Target: pod,
	}
	//deleteAct := PodUpdate{
	//	Action: DeleteAction,
	//	Target: pod,
	//}
	createMsg, err := json.Marshal(createAct)
	if err != nil {
		fmt.Println(err.Error())
	}

	//deleteMsg, err := json.Marshal(deleteAct)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	// after 5s, create the pod
	// after 1min, delete the pod
	go func() {
		createTimer := time.NewTimer(time.Second * 5)
		//deleteTimer := time.NewTimer(time.Minute)
		for i := 0; i < 2; i++ {
			select {
			//case <-deleteTimer.C:
			//	fmt.Println("Now delete the pod")
			//	listwatch.Publish(PodUpdateTopic, deleteMsg)
			case <-createTimer.C:
				fmt.Println("Now create the pod")
				listwatch.Publish(PodUpdateTopic, createMsg)
			}
		}
	}()

	kl := NewKubelet()
	kl.Run()
}
