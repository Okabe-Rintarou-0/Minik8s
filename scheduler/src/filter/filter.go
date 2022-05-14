package filter

import (
	"minik8s/apiObject"
	"minik8s/entity"
)

type Filter interface {
	Filter(pod *apiObject.Pod, nodes []*entity.NodeStatus) (filtered []*entity.NodeStatus)
}

func Default() Filter {
	return &filter{}
}

type filter struct{}

func (f *filter) containsAllLabels(nodeSelector, labels apiObject.Labels) bool {
	for key, value := range nodeSelector {
		if labelValue, exists := labels[key]; !exists || labelValue != value {
			return false
		}
	}
	return true
}

func (f *filter) Filter(pod *apiObject.Pod, nodes []*entity.NodeStatus) (filtered []*entity.NodeStatus) {
	nodeSelector := pod.NodeSelector()
	for _, node := range nodes {
		if f.containsAllLabels(nodeSelector, node.Labels) {
			filtered = append(filtered, node)
		}
	}
	return
}
