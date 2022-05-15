package etcd

import (
	"minik8s/apiObject"
	"minik8s/entity"
)

type ReplicaSetInfo struct {
	ReplicaSet apiObject.ReplicaSet
	Status     entity.ReplicaSetStatus
}
