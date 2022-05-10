package main

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/controller/controller"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util"
	"time"
)

func main() {
	content, err := ioutil.ReadFile("../test/rs.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	rs := apiObject.ReplicaSet{}
	_ = yaml.Unmarshal(content, &rs)
	rs.Metadata.UID = uuid.NewV4().String()
	topic := util.ReplicaSetUpdateTopic()

	content, err = ioutil.ReadFile("../test/rs2.yaml")
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
