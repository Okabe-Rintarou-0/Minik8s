package scheduler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/scheduler/src/filter"
	"minik8s/scheduler/src/selector"
	"testing"
	"time"
)

var testNode1 = &entity.NodeStatus{
	Hostname: "node1",
	Ip:       "0.0.0.0",
	Labels: map[string]string{
		"os": "linux",
	},
	Lifecycle:  entity.NodeReady,
	Error:      "",
	CpuPercent: 50,
	MemPercent: 30,
	NumPods:    1,
	SyncTime:   time.Now(),
}

var testNode2 = &entity.NodeStatus{
	Hostname: "node2",
	Ip:       "0.0.0.0",
	Labels: map[string]string{
		"os": "windows",
	},
	Lifecycle:  entity.NodeReady,
	Error:      "",
	CpuPercent: 30,
	MemPercent: 50,
	NumPods:    0,
	SyncTime:   time.Now(),
}

var testPod = &apiObject.Pod{
	Base: apiObject.Base{},
	Spec: apiObject.PodSpec{
		RestartPolicy: "Always",
		NodeSelector: map[string]string{
			"os": "linux",
		},
		Containers: nil,
		Volumes:    nil,
	},
}

func TestFilter(t *testing.T) {
	f := filter.Default()
	testNodes := []*entity.NodeStatus{testNode1, testNode2}
	filtered := f.Filter(testPod, testNodes)
	assert.Equal(t, 1, len(filtered))
	assert.Equal(t, "node1", filtered[0].Hostname)

	testPod.Spec.NodeSelector["os"] = "windows"
	filtered = f.Filter(testPod, testNodes)
	assert.Equal(t, 1, len(filtered))
	assert.Equal(t, "node2", filtered[0].Hostname)
}

func TestSelector(t *testing.T) {
	testNodes := []*entity.NodeStatus{testNode1, testNode2}

	// test mim cpu
	s := selector.DefaultFactory.NewSelector(selector.MinimumCpuUtility)
	selected := s.Select(testNodes)
	assert.Equal(t, "node2", selected.Hostname)

	s = selector.DefaultFactory.NewSelector(selector.MinimumMemoryUtility)
	selected = s.Select(testNodes)
	assert.Equal(t, "node1", selected.Hostname)

	s = selector.DefaultFactory.NewSelector(selector.MinimumNumPods)
	selected = s.Select(testNodes)
	assert.Equal(t, "node2", selected.Hostname)

	s = selector.DefaultFactory.NewSelector(selector.MaximumNumPods)
	selected = s.Select(testNodes)
	assert.Equal(t, "node1", selected.Hostname)

	s = selector.DefaultFactory.NewSelector(selector.Random)
	numOne, numTwo := 0, 0
	for i := 0; i < 1000; i++ {
		selected = s.Select(testNodes)
		if selected.Hostname == "node1" {
			numOne++
		} else {
			numTwo++
		}
	}
	fmt.Println(numOne, numTwo)
	bigger := float64(numOne)
	smaller := float64(numTwo)
	if numTwo > numOne {
		bigger = float64(numTwo)
		smaller = float64(numOne)
	}
	fmt.Println(smaller / bigger)
	assert.Less(t, 0.9, smaller/bigger)
}
