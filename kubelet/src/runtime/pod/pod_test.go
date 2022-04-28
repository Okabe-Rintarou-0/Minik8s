package pod

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
	pod := apiObject.Pod{}
	_ = yaml.Unmarshal(content, &pod)
	return &pod
}

func TestGetPodStatuses(t *testing.T) {
	var err error
	pod := readPod("./testPod.yaml")

	u1 := uuid.NewV4()
	fmt.Printf("UUID for test: %s\n", u1)
	pod.Metadata.UID = u1.String()

	pm := NewPodManager()
	err = pm.CreatePod(pod)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)

	pod = readPod("./testPod2.yaml")
	u2 := uuid.NewV4()
	fmt.Printf("UUID for test: %s\n", u2)
	pod.Metadata.UID = u2.String()

	err = pm.CreatePod(pod)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)

	podStatuses, err := pm.GetPodStatuses()
	assert.Nil(t, err)
	for _, ps := range podStatuses {
		fmt.Println(ps.ID)
	}
}

func TestCreatePod(t *testing.T) {
	var err error
	pod := readPod("./testPod.yaml")

	u1 := uuid.NewV4()
	fmt.Printf("UUID for test: %s\n", u1)
	pod.Metadata.UID = u1.String()

	pm := NewPodManager()
	err = pm.CreatePod(pod)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)

	podStatus, err := pm.GetPodStatus(pod)
	assert.Nil(t, err)
	for _, cs := range podStatus.ContainerStatuses {
		fmt.Println(cs.Name, cs.Name, cs.State, cs.CreatedAt.String())
	}
}

func TestStartContainer(t *testing.T) {
	var err error
	pod := readPod("./testPod.yaml")

	u1 := uuid.NewV4()
	fmt.Printf("UUID for test: %s\n", u1)
	pod.Metadata.UID = u1.String()

	pm := &podManager{
		cm: container.NewContainerManager(),
		im: image.NewImageManager(),
	}
	err = pm.startPauseContainer(pod)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = pm.startCommonContainer(pod, &pod.Spec.Containers[0])
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)

	containers, err := pm.cm.ListContainers(&container.ContainerListConfig{
		Quiet:  false,
		Size:   false,
		All:    true,
		Latest: false,
		Since:  "",
		Before: "",
		Limit:  5,
		LabelSelector: container.LabelSelector{
			KubernetesPodUIDLabel: pod.UID(),
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
