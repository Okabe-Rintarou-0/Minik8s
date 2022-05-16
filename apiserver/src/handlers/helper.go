package handlers

import (
	"encoding/json"
	"fmt"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"path"
	"strings"
)

func parseTargetName(targetName string) (namespace, name string) {
	parts := strings.Split(targetName, "/")
	numParts := len(parts)
	if numParts == 1 {
		return "default", targetName
	} else {
		return parts[0], strings.Join(parts[1:], "/")
	}
}

func existsTarget(target *apiObject.ScaleTargetRef) bool {
	etcdURL := path.Join(url.ReplicaSetURL, target.Namespace(), target.Name())
	rs, err := etcd.Get(etcdURL)
	if rs == "" || err != nil {
		return false
	}
	parseErr := json.Unmarshal([]byte(rs), &apiObject.ReplicaSet{})
	return parseErr == nil
}

func addHPA(hpa *apiObject.HorizontalPodAutoscaler) (err error) {
	if !existsTarget(hpa.Target()) {
		return fmt.Errorf("target %s/%s does not exits", hpa.TargetMetadata().Namespace, hpa.TargetMetadata().UID)
	}

	var hpaJson []byte
	if hpaJson, err = json.Marshal(hpa); err != nil {
		return
	}

	// exists?
	etcdURL := path.Join(url.HPAURL, hpa.Namespace(), hpa.Name())
	if hpaJsonStr, err := etcd.Get(etcdURL); err == nil {
		getHPA := &apiObject.HorizontalPodAutoscaler{}
		if err = json.Unmarshal([]byte(hpaJsonStr), getHPA); err == nil {
			return fmt.Errorf("hpa %s/%s already exists", getHPA.Namespace(), getHPA.Name())
		}
	}

	if err = etcd.Put(etcdURL, string(hpaJson)); err != nil {
		return
	}

	hpaUpdateMsg, _ := json.Marshal(entity.HPAUpdate{
		Action: entity.CreateAction,
		Target: *hpa,
	})

	listwatch.Publish(topicutil.HPAUpdateTopic(), hpaUpdateMsg)
	return nil
}
