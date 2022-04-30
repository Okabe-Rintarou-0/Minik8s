package kubelet

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"minik8s/apiObject"
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
	updates := make(chan PodUpdate, 5)
	pod := readPod("./testPod.yaml")
	pod.Metadata.UID = uuid.NewV4().String()
	updates <- PodUpdate{
		Action: CreateAction,
		Target: pod,
	}
	kl := NewKubelet()

	// after 15s, delete the pod
	go func() {
		timer := time.NewTimer(time.Second * 15)
		select {
		case <-timer.C:
			fmt.Println("Now delete the pod")
			updates <- PodUpdate{
				Action: DeleteAction,
				Target: pod,
			}
		}
	}()

	kl.Run(updates)
}

func TestParse(t *testing.T) {
	str := "/k8s_POD_test_test_87653ecc-04c3-4b32-b4e3-94346b968ede_0"
	succ, containerName, podName, namespace, uid, count := podutil.ParseContainerFullName(str)
	fmt.Println(succ, containerName, podName, namespace, uid, count)
}
