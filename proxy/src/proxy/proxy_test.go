package proxy

import (
	"encoding/json"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"minik8s/util/uidutil"
	"testing"
)

func TestPublish(t *testing.T) {
	update1 := entity.ServiceUpdate{}
	update1.Action = entity.CreateAction

	update1.Target.Service.Metadata.UID = uidutil.New()
	update1.Target.Service.Spec.ClusterIP = "10.96.1.1"
	update1.Target.Service.Spec.Ports = make([]apiObject.ServicePort, 1)
	update1.Target.Service.Spec.Ports[0] = apiObject.ServicePort{
		Name:       "test0",
		Port:       32220,
		TargetPort: 23333,
	}
	//update1.Target.Service.Spec.Ports[1] = apiObject.ServicePort{
	//	Name:       "test1",
	//	Port:       32221,
	//	TargetPort: 23333,
	//}
	//update1.Target.Service.Spec.Ports[2] = apiObject.ServicePort{
	//	Name:       "test1",
	//	Port:       32222,
	//	TargetPort: 23334,
	//}

	update1.Target.Endpoints = make([]apiObject.Endpoint, 2)
	update1.Target.Endpoints[0].UID = uidutil.New()
	update1.Target.Endpoints[0].IP = "127.0.0.1"
	update1.Target.Endpoints[0].Port = "23333"
	update1.Target.Endpoints[1].UID = uidutil.New()
	update1.Target.Endpoints[1].IP = "127.0.0.1"
	update1.Target.Endpoints[1].Port = "23334"
	msg, _ := json.Marshal(update1)
	listwatch.Publish(topicutil.ServiceUpdateTopic(), msg)
}
