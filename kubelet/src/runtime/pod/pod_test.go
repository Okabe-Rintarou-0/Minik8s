package pod

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testDocker/apiObject"
	"testDocker/kubelet/src/runtime/container"
	"testDocker/kubelet/src/runtime/image"
	"testing"
)

func TestStartContainer(t *testing.T) {
	var err error

	content, err := ioutil.ReadFile("./testPod.yaml")
	pod := apiObject.Pod{}
	if err == nil {
		err := yaml.Unmarshal(content, &pod)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("%v\n", pod)
	} else {
		fmt.Println(err.Error())
	}
	pm := &podManager{
		cm: container.NewContainerManager(),
		im: image.NewImageManager(),
	}
	err = pm.startPauseContainer(&pod)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = pm.startCommonContainer(&pod, &pod.Spec.Containers[0])
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
}
