package handlers

import (
	"encoding/json"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"minik8s/util/uidutil"
)

func addHPA(hpa *apiObject.HorizontalPodAutoscaler) error {
	hpa.Metadata.UID = uidutil.New()
	log("receive hpa[ID = %v]: %v", hpa.UID(), hpa)
	msg, err := json.Marshal(entity.HPAUpdate{
		Action: entity.CreateAction,
		Target: *hpa,
	})
	if err == nil {
		listwatch.Publish(topicutil.HPAUpdateTopic(), msg)
	}
	return err
}
