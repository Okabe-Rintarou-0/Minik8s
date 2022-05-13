package selector

import (
	"math/rand"
	"minik8s/entity"
)

type Selector interface {
	Select(nodes []*entity.NodeStatus) *entity.NodeStatus
}

type selector struct {
}

func New() Selector {
	return &selector{}
}

func (s *selector) Select(nodes []*entity.NodeStatus) *entity.NodeStatus {
	numNodes := len(nodes)
	if numNodes == 0 {
		return nil
	}
	randomIdx := 0
	if numNodes > 1 {
		randomIdx = rand.Intn(numNodes - 1)
	}
	return nodes[randomIdx]
}
