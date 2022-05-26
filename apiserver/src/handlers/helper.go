package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	hpaController "minik8s/controller/src/controller/hpa"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"minik8s/util/uidutil"
	"path"
	"strings"
	"time"
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

func getTarget(target *apiObject.ScaleTargetRef) *apiObject.ReplicaSet {
	etcdURL := path.Join(url.ReplicaSetURL, target.Namespace(), target.Name())
	raw, err := etcd.Get(etcdURL)
	if err != nil {
		return nil
	}
	rs := &apiObject.ReplicaSet{}
	if parseErr := json.Unmarshal([]byte(raw), rs); parseErr != nil {
		return nil
	}
	return rs
}

func addHPA(hpa *apiObject.HorizontalPodAutoscaler) (err error) {
	var rs *apiObject.ReplicaSet
	if rs = getTarget(hpa.Target()); rs == nil {
		return fmt.Errorf("target %s/%s does not exits", hpa.TargetMetadata().Namespace, hpa.TargetMetadata().Name)
	}
	hpa.SetTarget(rs)
	hpa.Metadata.UID = uidutil.New()
	if hpa.ScaleInterval() == 0 {
		hpa.Spec.ScaleInterval = hpaController.DefaultScaleInterval
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

	log("metrics %+v", hpa.Metrics())
	if err = etcd.Put(etcdURL, string(hpaJson)); err != nil {
		return
	}

	etcdHPAStatusURL := path.Join(url.HPAURL, "status", hpa.Namespace(), hpa.Name())
	var hpaStatusJson []byte
	if hpaStatusJson, err = json.Marshal(entity.HPAStatus{
		ID:        hpa.UID(),
		Name:      hpa.Name(),
		Namespace: hpa.Namespace(),
		Labels:    hpa.Labels(),
		Lifecycle: entity.HPACreated,
		Metrics:   "Unknown",
		Benchmark: 0,
		Error:     "",
		SyncTime:  time.Now(),
	}); err == nil {
		_ = etcd.Put(etcdHPAStatusURL, string(hpaStatusJson))
	}

	hpaUpdateMsg, _ := json.Marshal(entity.HPAUpdate{
		Action: entity.CreateAction,
		Target: *hpa,
	})

	listwatch.Publish(topicutil.HPAUpdateTopic(), hpaUpdateMsg)
	return nil
}

func getPutForm(body io.ReadCloser) (form map[string]string) {
	defer body.Close()
	content, _ := ioutil.ReadAll(body)
	if err := json.Unmarshal(content, &form); err != nil {
		return nil
	}
	return
}
