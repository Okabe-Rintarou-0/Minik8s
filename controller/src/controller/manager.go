package controller

import (
	"minik8s/controller/src/cache"
	"minik8s/controller/src/controller/hpa"
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
}

func (m *manager) Start() {
	m.cacheManager.Start()
	go m.replicaSetController.Run()
	go m.hpaController.Run()
	wait.Forever()
}

func NewControllerManager() Manager {
	m := &manager{}
	m.cacheManager = cache.NewManager()
	m.replicaSetController = replicaSet.NewController(m.cacheManager)
	m.hpaController = hpa.NewController(m.cacheManager)
	m.cacheManager.SetPodStatusUpdateHook(m.replicaSetController.Sync)
	return m
}
