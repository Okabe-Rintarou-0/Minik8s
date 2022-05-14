package cache

import (
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/httputil"
	"minik8s/util/logger"
	"time"
)

const (
	nodeStatusFullSyncPeriod = time.Minute
)

func (m *manager) syncNodeChange(cachedNodeStatuses, serverNodeStatuses []*entity.NodeStatus) (toDelete, toAdd []*entity.NodeStatus) {
	nodeStatusMap := map[string]*entity.NodeStatus{}
	for _, cached := range cachedNodeStatuses {
		nodeStatusMap[cached.Hostname] = cached
	}
	for _, server := range serverNodeStatuses {
		if _, exists := nodeStatusMap[server.Hostname]; exists {
			m.nodeStatusCache.Update(server.Hostname, server)
			delete(nodeStatusMap, server.Hostname)
		} else {
			toAdd = append(toAdd, server)
		}
	}

	for _, cached := range nodeStatusMap {
		toDelete = append(toDelete, cached)
	}

	return
}

func (m *manager) dealWithToAdd(toAdd []*entity.NodeStatus) {
	for _, ns := range toAdd {
		m.nodeStatusCache.Add(ns.Hostname, ns)
		//log("Add Node Status[host = %s]", ns.Hostname)
	}
}

func (m *manager) dealWithToDelete(toDelete []*entity.NodeStatus) {
	for _, ns := range toDelete {
		m.nodeStatusCache.Add(ns.Hostname, ns)
		//log("Add Node Status[host = %s]", ns.Hostname)
	}
}

func (m *manager) fullSyncNodeStatuses() {
	log("Full Sync NodeStatuses with api-server")
	// Step 1: Get all node statuses from api-server
	var serverNodeStatuses []*entity.NodeStatus
	err := httputil.GetAndUnmarshal(url.Prefix+url.NodeURL, &serverNodeStatuses)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// Step 2: Get the lock
	m.nodeStatusLock.Lock()
	defer m.nodeStatusLock.Unlock()

	// Step 3: compute node to delete
	cachedNodeStatuses := m.getNodeStatusesInternal()
	toDelete, toAdd := m.syncNodeChange(cachedNodeStatuses, serverNodeStatuses)

	// Step 4: Deal with toDelete & toAdd
	log("To Delete: %v", toDelete)
	log("To Add: %v", toAdd)
	m.dealWithToAdd(toAdd)
	m.dealWithToDelete(toDelete)
	///TODO implement it
}
