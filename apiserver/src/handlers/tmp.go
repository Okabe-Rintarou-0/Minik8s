package handlers

import (
	"minik8s/entity"
	"minik8s/util/uidutil"
	"time"
)

func podStatusForTest() *entity.PodStatus {
	return &entity.PodStatus{
		ID:        uidutil.New(),
		Name:      "example",
		Labels:    nil,
		Namespace: "default",
		Lifecycle: entity.PodContainerCreating,
		SyncTime:  time.Now(),
	}
}

func getPodDescriptionForTest(name string) *entity.PodDescription {
	var logs []entity.PodStatusLogEntry

	podStatus := podStatusForTest()
	logs = append(logs, entity.PodStatusLogEntry{
		Status: podStatus.Lifecycle,
		Time:   podStatus.SyncTime,
		Error:  podStatus.Error,
	})

	logs = append(logs, entity.PodStatusLogEntry{
		Status: entity.PodRunning,
		Time:   time.Now().Add(time.Minute * 30),
		Error:  "",
	})

	return &entity.PodDescription{
		CurrentStatus: *podStatusForTest(),
		Logs:          logs,
	}
}
