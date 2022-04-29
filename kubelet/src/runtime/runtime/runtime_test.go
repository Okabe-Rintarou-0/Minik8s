package runtime

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/kubelet/src/runtime/image"
	"testing"
)

func readPod(podPath string) *apiObject.Pod {
	content, _ := ioutil.ReadFile(podPath)
	testPod := apiObject.Pod{}
	_ = yaml.Unmarshal(content, &testPod)
	return &testPod
}

func TestGetPodStatuses(t *testing.T) {
	var err error
	testPod := readPod("./testPod.yaml")

	u1 := uuid.NewV4()
	fmt.Printf("UUID for test: %s\n", u1)
	testPod.Metadata.UID = u1.String()

	rm := NewPodManager()
	err = rm.CreatePod(testPod)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)

	testPod = readPod("./testPod2.yaml")
	u2 := uuid.NewV4()
	fmt.Printf("UUID for test: %s\n", u2)
	testPod.Metadata.UID = u2.String()

	err = rm.CreatePod(testPod)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)

	podStatuses, err := rm.GetPodStatuses()
	assert.Nil(t, err)
	for _, ps := range podStatuses {
		fmt.Println(ps.ID)
	}
}

func TestCreatePod(t *testing.T) {
	var err error
	testPod := readPod("./testPod.yaml")

	u1 := uuid.NewV4()
	fmt.Printf("UUID for test: %s\n", u1)
	testPod.Metadata.UID = u1.String()

	rm := NewPodManager()
	err = rm.CreatePod(testPod)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)

	podStatus, err := rm.GetPodStatus(testPod)
	assert.Nil(t, err)
	for _, cs := range podStatus.ContainerStatuses {
		fmt.Println(cs.Name, cs.Name, cs.State, cs.CreatedAt.String())
	}
}

func TestStartContainer(t *testing.T) {
	var err error
	testPod := readPod("./testPod.yaml")

	u1 := uuid.NewV4()
	fmt.Printf("UUID for test: %s\n", u1)
	testPod.Metadata.UID = u1.String()

	rm := &runtimeManager{
		cm: container.NewContainerManager(),
		im: image.NewImageManager(),
	}
	err = rm.startPauseContainer(testPod)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = rm.startCommonContainer(testPod, &testPod.Spec.Containers[0])
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)

	containers, err := rm.cm.ListContainers(&container.ContainerListConfig{
		Quiet:  false,
		Size:   false,
		All:    true,
		Latest: false,
		Since:  "",
		Before: "",
		Limit:  5,
		LabelSelector: container.LabelSelector{
			KubernetesPodUIDLabel: testPod.UID(),
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
	for _, c := range containers {
		fmt.Println("Got container", c.Name, c.ID, c.Image)
	}
}
