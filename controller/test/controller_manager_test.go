package test

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/controller/src/controller"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"minik8s/util/uidutil"
	"testing"
	"time"
)

func TestControllerManager(t *testing.T) {
	content, err := ioutil.ReadFile("./rs.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	rs := apiObject.ReplicaSet{}
	_ = yaml.Unmarshal(content, &rs)
	rs.Metadata.UID = uidutil.New()
	topic := topicutil.ReplicaSetUpdateTopic()

	content, err = ioutil.ReadFile("./rs2.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	rs2 := apiObject.ReplicaSet{}
	_ = yaml.Unmarshal(content, &rs2)
	rs2.Metadata.UID = rs.UID()

	go func() {
		<-time.Tick(time.Second * 5)
		fmt.Println("Create a rs")
		msg, _ := json.Marshal(entity.ReplicaSetUpdate{
			Action: entity.CreateAction,
			Target: rs,
		})
		listwatch.Publish(topic, msg)

		<-time.Tick(time.Second * 25)
		fmt.Println("Update a rs")
		msg, _ = json.Marshal(entity.ReplicaSetUpdate{
			Action: entity.UpdateAction,
			Target: rs2,
		})
		listwatch.Publish(topic, msg)

		<-time.Tick(time.Minute)
		fmt.Println("Delete the rs")
		msg, _ = json.Marshal(entity.ReplicaSetUpdate{
			Action: entity.DeleteAction,
			Target: rs,
		})
		listwatch.Publish(topic, msg)
	}()

	cm := controller.NewControllerManager()
	cm.Start()
}
