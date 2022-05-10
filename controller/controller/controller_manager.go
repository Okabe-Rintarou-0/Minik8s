package controller

import (
	"minik8s/controller/controller/cache"
	"minik8s/controller/controller/replicaSet"
)

type Manager interface {
	Start()
}

type manager struct {
	cacheManager         cache.Manager
	replicaSetController replicaSet.Controller
}

func (m *manager) Start() {
	go m.cacheManager.Start()
	m.replicaSetController.Run()
}

func NewControllerManager() Manager {
	m := &manager{}
	m.cacheManager = cache.NewManager()
	m.replicaSetController = replicaSet.NewController(m.cacheManager)
	m.cacheManager.SetPodStatusUpdateHook(m.replicaSetController.Sync)
	return m
}
