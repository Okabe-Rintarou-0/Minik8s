package etcd

import (
	"minik8s/apiObject"
	"minik8s/entity"
)

type PodInfo struct {
	Pod    apiObject.Pod
	Status entity.PodStatus
}
