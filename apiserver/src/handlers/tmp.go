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
