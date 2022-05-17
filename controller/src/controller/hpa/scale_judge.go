package hpa

import (
	"fmt"
	"math"
	"minik8s/entity"
	"minik8s/util/mathutil"
)

type ScaleJudge interface {
	// Judge returns the number of replicas according to given status and metrics
	Judge(status *entity.ReplicaSetStatus) int
	Metrics() entity.HPAMetrics
	Benchmark() float64
}

func NewCpuScaleJudge(benchmark float64, minReplicas int, maxReplicas int) ScaleJudge {
	return &cpuScaleJudge{
		benchmark:   benchmark,
		minReplicas: minReplicas,
		maxReplicas: maxReplicas,
	}
}

func NewMemoryScaleJudge(benchmark float64, minReplicas int, maxReplicas int) ScaleJudge {
	return &memScaleJudge{
		benchmark:   benchmark,
		minReplicas: minReplicas,
		maxReplicas: maxReplicas,
	}
}

// FakeScaleJudge is for test only, do not use it!🥰
func FakeScaleJudge() ScaleJudge {
	return &fakeScaleJudge{}
}

type cpuScaleJudge struct {
	benchmark   float64
	minReplicas int
	maxReplicas int
}

func (c *cpuScaleJudge) Metrics() entity.HPAMetrics {
	return entity.CpuMetrics
}

func (c *cpuScaleJudge) Benchmark() float64 {
	return c.benchmark
}

func (c *cpuScaleJudge) Judge(status *entity.ReplicaSetStatus) int {
	cpuPercent := status.CpuPercent
	ratio := c.benchmark / cpuPercent
	numReplicas := mathutil.Clamp(int(math.Round(ratio*float64(status.NumReplicas))), c.minReplicas, c.maxReplicas)
	fmt.Printf("[CPU judge] Benchmark = %v, cpuPercent = %v, So num replicas should be: %d\n", c.benchmark, cpuPercent, numReplicas)
	return numReplicas
}

type memScaleJudge struct {
	benchmark   float64
	minReplicas int
	maxReplicas int
}

func (m *memScaleJudge) Metrics() entity.HPAMetrics {
	return entity.MemoryMetrics
}

func (m *memScaleJudge) Benchmark() float64 {
	return m.benchmark
}

func (m *memScaleJudge) Judge(status *entity.ReplicaSetStatus) int {
	memPercent := status.MemPercent
	ratio := memPercent / m.benchmark
	numReplicas := mathutil.Clamp(int(math.Round(ratio*float64(status.NumReplicas))), m.minReplicas, m.maxReplicas)
	fmt.Printf("[Mem judge] Num replicas should be: %d\n", numReplicas)
	return numReplicas
}

type fakeScaleJudge struct {
}

func (f *fakeScaleJudge) Metrics() entity.HPAMetrics {
	return entity.UnknownMetrics
}

func (f *fakeScaleJudge) Benchmark() float64 {
	return 0
}

func (f *fakeScaleJudge) Judge(status *entity.ReplicaSetStatus) int {
	return 2
}
