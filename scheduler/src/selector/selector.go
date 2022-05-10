package selector

import (
	"math/rand"
	"minik8s/apiObject"
)

type Selector interface {
	Select(nodes []*apiObject.Node) *apiObject.Node
}

type selector struct {
}

func New() Selector {
	return &selector{}
}

func (s *selector) Select(nodes []*apiObject.Node) *apiObject.Node {
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
