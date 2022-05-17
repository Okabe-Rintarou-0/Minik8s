package hpa

import (
	"context"
	"encoding/json"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/controller/src/cache"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/httputil"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"path"
	"strconv"
	"time"
)

var logWorker = logger.Log("HPA worker")

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

// getReplicaSetStatus get rs status from cache
func (w *worker) getReplicaSetStatus() *entity.ReplicaSetStatus {
	targetReplicaSet := w.target.Target()
	UID := targetReplicaSet.UID()
	return w.cacheManager.GetReplicaSetStatus(UID)
}

func (w *worker) updateReplicaSetToApiServer(namespace, name string, numReplicas int) {
	URL := url.Prefix + path.Join(url.ReplicaSetURL, namespace, name)
	resp := httputil.PutForm(URL, map[string]string{
		"replicas": strconv.Itoa(numReplicas),
	})
	logWorker("update rs and get resp: %s", resp)
}

func (w *worker) updateReplicaSet(namespace, name string, numReplicas int) {
	w.updateReplicaSetToApiServer(namespace, name, numReplicas)
}

func (w *worker) publishHPAStatus(targetReplicas int, targetStatus *entity.ReplicaSetStatus) {
	hpaStatus := &entity.HPAStatus{
		ID:          w.target.UID(),
		Name:        w.target.Name(),
		Namespace:   w.target.Namespace(),
		Labels:      w.target.Labels(),
		Lifecycle:   entity.HPAReady,
		MinReplicas: w.target.MinReplicas(),
		MaxReplicas: w.target.MaxReplicas(),
		Metrics:     w.scaleJudge.Metrics(),
		Benchmark:   w.scaleJudge.Benchmark(),
		NumReady:    targetStatus.NumReady,
		NumTarget:   targetReplicas,
		Error:       "",
		SyncTime:    time.Now(),
	}

	hpaStatusJson, _ := json.Marshal(hpaStatus)
	logWorker("Publish hpaStatus %v[ID = %s]", hpaStatus.Lifecycle.String(), hpaStatus.ID)
	listwatch.Publish(topicutil.HPAStatusTopic(), hpaStatusJson)
}

func (w *worker) syncLoopIteration() bool {
	// Step 1: Get replicaSet status from Cache
	replicaSetStatus := w.getReplicaSetStatus()
	logWorker("Got replicaSet status: %v", replicaSetStatus)
	if replicaSetStatus == nil {
		logWorker("Something bad happens...")
		return true
	}

	// judge the numReplicas we need
	numReplicas := w.scaleJudge.Judge(replicaSetStatus)
	logWorker("Judge result is %d num replicas", numReplicas)

	diff := replicaSetStatus.NumReplicas - numReplicas
	logWorker("Auto scale result: diff = %d", diff)
	if diff > 0 {
		go w.updateReplicaSet(replicaSetStatus.Namespace, replicaSetStatus.Name, replicaSetStatus.NumReplicas-1)
	} else if diff < 0 {
		go w.updateReplicaSet(replicaSetStatus.Namespace, replicaSetStatus.Name, replicaSetStatus.NumReplicas+1)
	}
	w.publishHPAStatus(numReplicas, replicaSetStatus)
	return true
}

func (w *worker) syncLoop() {
	period := time.Second * time.Duration(w.target.ScaleInterval())
	logWorker("Scale interval is %d", w.target.ScaleInterval())
	tick := time.Tick(period)
	for w.syncLoopIteration() {
		select {
		case <-tick:
			continue
		case <-w.ctx.Done():
			logWorker("Stop working on hpa[ID = %s]\n", w.target.UID())
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
	//return FakeScaleJudge()

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
