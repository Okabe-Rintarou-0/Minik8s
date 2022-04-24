package pod

import (
	"testDocker/src/runtime/container"
	"testDocker/src/types"
)

type Pod struct {
	ID         types.UID
	Name       string
	Namespace  string
	Containers []*container.Container
}
