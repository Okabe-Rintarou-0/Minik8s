package test

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
	rs.Metadata.UID = uuid.NewV4().String()
	topic := util.ReplicaSetUpdateTopic()

	go func() {
		<-time.Tick(time.Second * 5)
		fmt.Println("create a rs")
		msg, _ := json.Marshal(entity.ReplicaSetUpdate{
			Action: entity.CreateAction,
			Target: rs,
		})
		listwatch.Publish(topic, msg)
	}()

	cm := controller.NewControllerManager()
	cm.Start()
}
