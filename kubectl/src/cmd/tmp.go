package cmd

import (
	uuid "github.com/satori/go.uuid"
	"minik8s/entity"
	"time"
)

func podStatusForTest() *entity.PodStatus {
	return &entity.PodStatus{
		ID:        uuid.NewV4().String(),
		Name:      "example",
		Labels:    nil,
		Namespace: "default",
		Status:    entity.ContainerCreating,
		SyncTime:  time.Now(),
	}
}

func podDescriptionForTest(name string) *entity.PodDescription {
	//TODO just for test now, replace it with api-server

	var logs []entity.PodStatusLogEntry

	podStatus := podStatusForTest()
	logs = append(logs, entity.PodStatusLogEntry{
		Status: podStatus.Status,
		Time:   podStatus.SyncTime,
		Error:  podStatus.Error,
	})

	logs = append(logs, entity.PodStatusLogEntry{
		Status: entity.Running,
		Time:   time.Now().Add(time.Minute * 30),
		Error:  "",
	})

	return &entity.PodDescription{
		CurrentStatus: *podStatusForTest(),
		Logs:          logs,
	}
}
