package hpa

import (
	"context"
	"encoding/json"
	"fmt"
	"minik8s/apiObject"
	"minik8s/controller/src/cache"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"time"
)

const syncPeriodSeconds = 30

type Worker interface {
	Run()
	SetTarget(rs *apiObject.HorizontalPodAutoscaler)
}

type worker struct {
	cacheManager cache.Manager
	target       *apiObject.HorizontalPodAutoscaler
	scaleJudge   ScaleJudge
	ctx          context.Context
}

func (w *worker) SetTarget(hpa *apiObject.HorizontalPodAutoscaler) {
	if hpa != nil {
		w.target = hpa
	}
}

// testMap is just for *TEST*, do not use it.
// We have to save pod into map, because we don't have an api-server now
// All we can do is to *Mock*
var testMap = map[string]*apiObject.ReplicaSet{}

// AddRsForTest TODO just for test
func AddRsForTest(rs *apiObject.ReplicaSet) {
	testMap[rs.FullName()] = rs
}

// getReplicaSetStatus get rs status from cache
func (w *worker) getReplicaSetStatus() *entity.ReplicaSetStatus {
	targetReplicaSet := w.target.Target()
	return w.cacheManager.GetReplicaSetStatus(targetReplicaSet.FullName())
}

func (w *worker) updateReplicaSetToApiServerForTest(fullName string, numReplicas int) {
	rs := testMap[fullName]
	rs.SetReplicas(numReplicas)
	fmt.Println("[HPA]Update rs test pod", rs)
	msg, _ := json.Marshal(entity.ReplicaSetUpdate{
		Action: entity.UpdateAction,
		Target: *rs,
	})
	topic := topicutil.ReplicaSetUpdateTopic()
	listwatch.Publish(topic, msg)
}

func (w *worker) updateReplicaSet(fullName string, numReplicas int) {
	// Add two, so we can test the case that the number of existent pods is more than replicas
	// just for test now
	w.updateReplicaSetToApiServerForTest(fullName, numReplicas)
}

func (w *worker) syncLoopIteration() bool {
	// Step 1: Get replicaSet status from Cache
	replicaSetStatus := w.getReplicaSetStatus()
	fmt.Println("[HPA]Got replicaSet status:", replicaSetStatus)
	if replicaSetStatus == nil {
		fmt.Println("Something bad happens...🙄")
		return true
	}

	// judge the numReplicas we need
	numReplicas := w.scaleJudge.Judge(replicaSetStatus)
	fmt.Printf("[HPA]judge result is %d num replicas\n", numReplicas)

	diff := replicaSetStatus.NumReplicas - numReplicas
	fmt.Printf("[HPA]auto scale result: diff = %d\n", diff)
	if diff != 0 {
		go w.updateReplicaSet(replicaSetStatus.FullName(), numReplicas)
	}
	return true
}

func (w *worker) syncLoop() {
	tick := time.Tick(time.Second * syncPeriodSeconds)
	for {
		select {
		case <-tick:
			if !w.syncLoopIteration() {
				return
			}
		case <-w.ctx.Done():
			fmt.Printf("[HPA worker]stop working on hpa[ID = %s]", w.target.UID())
			return
		}
	}
}

func (w *worker) Run() {
	w.syncLoop()
}

func decideMetrics(hpa *apiObject.HorizontalPodAutoscaler) ScaleJudge {
	var (
		metrics     = hpa.Metrics()
		cpuPercent  = metrics.CPUUtilizationPercentage
		memPercent  = metrics.MemUtilizationPercentage
		minReplicas = hpa.MinReplicas()
		maxReplicas = hpa.MaxReplicas()
	)

	// for test only, delete it afterwards
	return FakeScaleJudge()

	// Cpu takes the highest priority
	if cpuPercent > 0 {
		return NewCpuScaleJudge(cpuPercent, minReplicas, maxReplicas)
	}

	// Then, memory
	if memPercent > 0 {
		return NewMemoryScaleJudge(memPercent, minReplicas, maxReplicas)
	}

	// What if both of them are not valid?
	// No worry, kubectl and api-server should judge its validity
	return nil
}

func NewWorker(ctx context.Context, target *apiObject.HorizontalPodAutoscaler, cacheManager cache.Manager) Worker {
	w := &worker{
		target:       target,
		cacheManager: cacheManager,
		ctx:          ctx,
	}
	w.scaleJudge = decideMetrics(target)
	return w
}
