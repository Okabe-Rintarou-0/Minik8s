package controller

import (
	"minik8s/controller/src/cache"
	"minik8s/controller/src/controller/gpu"
	"minik8s/controller/src/controller/hpa"
	"minik8s/controller/src/controller/node"
	"minik8s/controller/src/controller/replicaSet"
	"minik8s/util/wait"
)

type Manager interface {
	Start()
}

type manager struct {
	cacheManager         cache.Manager
	hpaController        hpa.Controller
	replicaSetController replicaSet.Controller
	nodeController       node.Controller
	gpuController        gpu.Controller
}

func (m *manager) Start() {
	m.cacheManager.Start()
	go m.replicaSetController.Run()
	go m.hpaController.Run()
	go m.nodeController.Run()
	go m.gpuController.Run()
	wait.Forever()
}

func NewControllerManager() Manager {
	m := &manager{}
	m.cacheManager = cache.NewManager()
	m.replicaSetController = replicaSet.NewController(m.cacheManager)
	m.hpaController = hpa.NewController(m.cacheManager)
	m.nodeController = node.NewController(m.cacheManager)
	m.gpuController = gpu.NewController()

	m.cacheManager.SetPodStatusUpdateHook(m.replicaSetController.Sync)
	m.cacheManager.SetReplicaSetFullSyncAddHook(m.replicaSetController.AddReplicaSet)
	m.cacheManager.SetReplicaSetFullSyncDeleteHook(m.replicaSetController.DeleteReplicaSet)
	m.cacheManager.SetHPAFullSyncAddHook(m.hpaController.AddHpa)
	m.cacheManager.SetHPAFullSyncDeleteHook(m.hpaController.DeleteHpa)
	return m
}
