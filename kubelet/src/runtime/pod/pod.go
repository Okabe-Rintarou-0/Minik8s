package pod

import (
	"testDocker/kubelet/src/runtime/container"
	"testDocker/kubelet/src/types"
)

type Pod struct {
	ID         types.UID
	Name       string
	Namespace  string
	Containers []*container.Container
}

type PodManager interface {
}
