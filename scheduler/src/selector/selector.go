package selector

import (
	"math"
	"math/rand"
	"minik8s/entity"
)

type Selector interface {
	Select(nodes []*entity.NodeStatus) *entity.NodeStatus
}

type randomSelector struct{}

func random() Selector {
	return &randomSelector{}
}

func (s *randomSelector) Select(nodes []*entity.NodeStatus) *entity.NodeStatus {
	numNodes := len(nodes)
	if numNodes == 0 {
		return nil
	}
	randomIdx := rand.Intn(numNodes)
	return nodes[randomIdx]
}

func minimumCpuUtility() Selector {
	return &minimumCpuUtilitySelector{}
}

type minimumCpuUtilitySelector struct{}

func (s *minimumCpuUtilitySelector) Select(nodes []*entity.NodeStatus) *entity.NodeStatus {
	if len(nodes) == 0 {
		return nil
	}

	minCpu := math.MaxFloat64
	minIdx := -1

	for i, node := range nodes {
		if node.CpuPercent < minCpu {
			minCpu = node.CpuPercent
			minIdx = i
		}
	}

	return nodes[minIdx]
}

func minimumMemoryUtility() Selector {
	return &minimumMemoryUtilitySelector{}
}

type minimumMemoryUtilitySelector struct{}

func (s *minimumNumPodsSelector) Select(nodes []*entity.NodeStatus) *entity.NodeStatus {
	if len(nodes) == 0 {
		return nil
	}

	minNumPods := math.MaxInt32
	minIdx := -1

	for i, node := range nodes {
		if node.NumPods < minNumPods {
			minNumPods = node.NumPods
			minIdx = i
		}
	}

	return nodes[minIdx]
}

func minimumNumPods() Selector {
	return &minimumNumPodsSelector{}
}

type minimumNumPodsSelector struct{}

func (s *minimumMemoryUtilitySelector) Select(nodes []*entity.NodeStatus) *entity.NodeStatus {
	if len(nodes) == 0 {
		return nil
	}

	minMem := math.MaxFloat64
	minIdx := -1

	for i, node := range nodes {
		if node.MemPercent < minMem {
			minMem = node.MemPercent
			minIdx = i
		}
	}

	return nodes[minIdx]
}

func maximumNumPods() Selector {
	return &maximumNumPodsSelector{}
}

type maximumNumPodsSelector struct{}

func (s *maximumNumPodsSelector) Select(nodes []*entity.NodeStatus) *entity.NodeStatus {
	if len(nodes) == 0 {
		return nil
	}

	maxNumPods := -1
	maxIdx := -1

	for i, node := range nodes {
		if node.NumPods > maxNumPods {
			maxNumPods = node.NumPods
			maxIdx = i
		}
	}

	return nodes[maxIdx]
}
