package cache

import (
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/httputil"
	"minik8s/util/logger"
	"path"
	"time"
)

const (
	nodeStatusFullSyncPeriod       = time.Minute
	podStatusFullSyncPeriod        = time.Second * 30
	replicaSetStatusFullSyncPeriod = time.Second * 50
	hpaStatusFullSyncPeriod        = time.Second * 50
)

func (m *manager) syncNodeChange(cachedNodeStatuses, serverNodeStatuses []*entity.NodeStatus) (toDelete, toAdd []*entity.NodeStatus) {
	nodeStatusMap := map[string]*entity.NodeStatus{}

	for _, cached := range cachedNodeStatuses {
		fullName := path.Join(cached.Namespace, cached.Hostname)
		nodeStatusMap[fullName] = cached
	}

	for _, server := range serverNodeStatuses {
		fullName := path.Join(server.Namespace, server.Hostname)
		if _, exists := nodeStatusMap[fullName]; exists {
			m.nodeStatusCache.Update(fullName, server)
			delete(nodeStatusMap, fullName)
		} else {
			toAdd = append(toAdd, server)
		}
	}

	for _, cached := range nodeStatusMap {
		toDelete = append(toDelete, cached)
	}
	return
}

func (m *manager) handleNodeAdd(toAdd []*entity.NodeStatus) {
	for _, ns := range toAdd {
		m.nodeStatusCache.Add(ns.Hostname, ns)
		//log("Add Node Status[host = %s]", ns.Hostname)
	}
}

func (m *manager) handleNodeDelete(toDelete []*entity.NodeStatus) {
	for _, ns := range toDelete {
		m.nodeStatusCache.Delete(ns.Hostname)
		//log("Add Node Status[host = %s]", ns.Hostname)
	}
}

func (m *manager) syncPodChange(cachedPodStatuses, serverPodStatuses []*entity.PodStatus) (toDelete, toAdd []*entity.PodStatus) {
	podStatusMap := map[string]*entity.PodStatus{}
	for _, cached := range cachedPodStatuses {
		podStatusMap[cached.ID] = cached
	}
	for _, server := range serverPodStatuses {
		if _, exists := podStatusMap[server.ID]; exists {
			m.podStatusCache.Update(server.ID, server)
			delete(podStatusMap, server.ID)
		} else {
			toAdd = append(toAdd, server)
		}
	}

	for _, cached := range podStatusMap {
		toDelete = append(toDelete, cached)
	}
	return
}

func (m *manager) handlePodAdd(toAdd []*entity.PodStatus) {
	for _, ps := range toAdd {
		m.podStatusCache.Add(ps.ID, ps)
		//log("Add Node Status[host = %s]", ns.Hostname)
	}
}

func (m *manager) handlePodDelete(toDelete []*entity.PodStatus) {
	for _, ps := range toDelete {
		m.podStatusCache.Delete(ps.ID)
		//log("Add Node Status[host = %s]", ns.Hostname)
	}
}

func (m *manager) syncReplicaSetChange(cachedReplicaSetStatuses, serverReplicaSetStatuses []*entity.ReplicaSetStatus) (toDelete, toAdd []*entity.ReplicaSetStatus) {
	replicaSetStatusMap := map[string]*entity.ReplicaSetStatus{}
	for _, cached := range cachedReplicaSetStatuses {
		replicaSetStatusMap[cached.ID] = cached
	}
	for _, server := range serverReplicaSetStatuses {
		if _, exists := replicaSetStatusMap[server.ID]; exists {
			m.replicaSetStatusCache.Update(server.ID, server)
			delete(replicaSetStatusMap, server.ID)
		} else {
			toAdd = append(toAdd, server)
		}
	}

	for _, cached := range replicaSetStatusMap {
		toDelete = append(toDelete, cached)
	}
	return
}

func (m *manager) syncHPAChange(cachedHPAStatuses, serverHPAStatuses []*entity.HPAStatus) (toDelete, toAdd []*entity.HPAStatus) {
	hpaStatusMap := map[string]*entity.HPAStatus{}
	for _, cached := range cachedHPAStatuses {
		hpaStatusMap[cached.ID] = cached
	}
	for _, server := range serverHPAStatuses {
		if _, exists := hpaStatusMap[server.ID]; exists {
			m.hpaStatusCache.Update(server.ID, server)
			delete(hpaStatusMap, server.ID)
		} else {
			toAdd = append(toAdd, server)
		}
	}

	for _, cached := range hpaStatusMap {
		toDelete = append(toDelete, cached)
	}
	return
}

func (m *manager) handleReplicaSetAdd(toAdd []*entity.ReplicaSetStatus) {
	for _, rs := range toAdd {
		m.replicaSetStatusCache.Add(rs.ID, rs)
		replicaSet := &apiObject.ReplicaSet{}
		URL := url.Prefix + path.Join(url.ReplicaSetURL, rs.Namespace, rs.Name)
		if err := httputil.GetAndUnmarshal(URL, replicaSet); err == nil {
			m.replicaSetFullSyncAddHook(replicaSet)
		}
		//log("Add Node Status[host = %s]", ns.Hostname)
	}
}

func (m *manager) handleReplicaSetDelete(toDelete []*entity.ReplicaSetStatus) {
	for _, rs := range toDelete {
		m.replicaSetStatusCache.Delete(rs.ID)
		replicaSet := &apiObject.ReplicaSet{}
		URL := url.Prefix + path.Join(url.ReplicaSetURL, rs.Namespace, rs.Name)
		if err := httputil.GetAndUnmarshal(URL, replicaSet); err == nil {
			m.replicaSetFullSyncDeleteHook(replicaSet)
		}
		//log("Add Node Status[host = %s]", ns.Hostname)
	}
}

func (m *manager) handleHPAAdd(toAdd []*entity.HPAStatus) {
	for _, hpas := range toAdd {
		m.hpaStatusCache.Add(hpas.ID, hpas)
		hpa := &apiObject.HorizontalPodAutoscaler{}
		URL := url.Prefix + path.Join(url.HPAURL, hpas.Namespace, hpas.Name)
		if err := httputil.GetAndUnmarshal(URL, hpa); err == nil {
			m.hpaFullSyncAddHook(hpa)
		}
		//log("Add Node Status[host = %s]", ns.Hostname)
	}
}

func (m *manager) handleHPADelete(toDelete []*entity.HPAStatus) {
	for _, hpas := range toDelete {
		m.hpaStatusCache.Delete(hpas.ID)
		hpa := &apiObject.HorizontalPodAutoscaler{}
		URL := url.Prefix + path.Join(url.HPAURL, hpas.Namespace, hpas.Name)
		if err := httputil.GetAndUnmarshal(URL, hpa); err == nil {
			m.hpaFullSyncDeleteHook(hpa)
		}
		//log("Add Delete Status[host = %s]", ns.Hostname)
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

	// Step 3: compute node to add or to delete
	cachedNodeStatuses := m.getNodeStatusesInternal()
	toDelete, toAdd := m.syncNodeChange(cachedNodeStatuses, serverNodeStatuses)

	// Step 4: Handle toDelete & toAdd
	//log("To Delete: %v", toDelete)
	//log("To Add: %v", toAdd)
	m.handleNodeAdd(toAdd)
	m.handleNodeDelete(toDelete)
}

func (m *manager) fullSyncPodStatuses() {
	log("Full Sync PodStatuses with api-server")
	// Step 1: Get all pod statuses from api-server
	var serverPodStatuses []*entity.PodStatus
	err := httputil.GetAndUnmarshal(url.Prefix+url.PodURL, &serverPodStatuses)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// Step 2: compute pod to add or to delete
	cachedPodStatuses := m.getPodStatusesInternal()
	toDelete, toAdd := m.syncPodChange(cachedPodStatuses, serverPodStatuses)

	// Step 3: Handle toDelete & toAdd
	//log("To Delete: %v", toDelete)
	//log("To Add: %v", toAdd)
	m.handlePodAdd(toAdd)
	m.handlePodDelete(toDelete)
}

func (m *manager) fullSyncReplicaSetStatuses() {
	log("Full Sync ReplicaSetStatuses with api-server")
	// Step 1: Get all replicaSet statuses from api-server
	var serverReplicaSetStatuses []*entity.ReplicaSetStatus
	err := httputil.GetAndUnmarshal(url.Prefix+url.ReplicaSetURL, &serverReplicaSetStatuses)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// Step 2: compute replicaset to add or to delete
	cachedReplicaSetStatuses := m.getReplicaSetStatusesInternal()
	for _, serverStatus := range serverReplicaSetStatuses {
		log("received from server: %+v", serverStatus)
	}
	toDelete, toAdd := m.syncReplicaSetChange(cachedReplicaSetStatuses, serverReplicaSetStatuses)

	// Step 3: Handle toDelete & toAdd
	log("To Delete: %v", toDelete)
	log("To Add: %v", toAdd)
	m.handleReplicaSetAdd(toAdd)
	m.handleReplicaSetDelete(toDelete)
}

func (m *manager) fullSyncHPAStatuses() {
	log("Full Sync HPAStatuses with api-server")
	// Step 1: Get all hpa statuses from api-server
	var serverHPAStatuses []*entity.HPAStatus
	err := httputil.GetAndUnmarshal(url.Prefix+url.HPAURL, &serverHPAStatuses)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// Step 2: compute replicaset to add or to delete
	cachedReplicaSetStatuses := m.getHPAStatusesInternal()
	toDelete, toAdd := m.syncHPAChange(cachedReplicaSetStatuses, serverHPAStatuses)

	// Step 3: Handle toDelete & toAdd
	//log("To Delete: %v", toDelete)
	//log("To Add: %v", toAdd)
	m.handleHPAAdd(toAdd)
	m.handleHPADelete(toDelete)
}
