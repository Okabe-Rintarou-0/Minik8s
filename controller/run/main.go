package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/controller/src/controller"
	"minik8s/util/parseutil"
	"minik8s/util/uidutil"
)

func main() {
	content, err := ioutil.ReadFile("../test/rs.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	rs := apiObject.ReplicaSet{}
	_ = yaml.Unmarshal(content, &rs)
	rs.Metadata.UID = uidutil.New()
	//topic := topicutil.ReplicaSetUpdateTopic()

	content, err = ioutil.ReadFile("../test/rs2.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	rs2 := apiObject.ReplicaSet{}
	_ = yaml.Unmarshal(content, &rs2)
	rs2.Metadata.UID = uidutil.New()

	content, _ = ioutil.ReadFile("../test/hpa.yaml")
	hpa, _ := parseutil.ParseHPA(content)
	hpa.Metadata.UID = uidutil.New()
	fmt.Printf("Read hpa: [min = %v, max = %v]\n", hpa.MinReplicas(), hpa.MaxReplicas())

	content, _ = ioutil.ReadFile("../test/hpa2.yaml")
	hpa2, _ := parseutil.ParseHPA(content)
	hpa2.Metadata.UID = uidutil.New()
	fmt.Printf("Read hpa2: [min = %v, max = %v]\n", hpa.MinReplicas(), hpa.MaxReplicas())

	//go func() {
	//<-time.Tick(time.Second * 5)
	//fmt.Println("Create rs")
	//msg, _ := json.Marshal(entity.ReplicaSetUpdate{
	//	Action: entity.CreateAction,
	//	Target: rs,
	//})
	//listwatch.Publish(topic, msg)
	//
	//<-time.Tick(time.Second * 5)
	//fmt.Println("Create rs2")
	//msg, _ = json.Marshal(entity.ReplicaSetUpdate{
	//	Action: entity.CreateAction,
	//	Target: rs2,
	//})
	//listwatch.Publish(topic, msg)
	//
	//<-time.Tick(time.Second * 15)
	//fmt.Println("Add a hpa")
	//msg, _ = json.Marshal(entity.HPAUpdate{
	//	Action: entity.CreateAction,
	//	Target: *hpa,
	//})
	//listwatch.Publish(topicutil.HPAUpdateTopic(), msg)
	//
	//<-time.Tick(time.Second * 5)
	//fmt.Println("Add hpa2")
	//msg, _ = json.Marshal(entity.HPAUpdate{
	//	Action: entity.CreateAction,
	//	Target: *hpa2,
	//})
	//listwatch.Publish(topicutil.HPAUpdateTopic(), msg)

	//<-time.Tick(time.Second * 25)
	//fmt.Println("Delete a hpa")
	//msg, _ = json.Marshal(entity.HPAUpdate{
	//	Action: entity.DeleteAction,
	//	Target: *hpa,
	//})
	//listwatch.Publish(topicutil.HPAUpdateTopic(), msg)

	//<-time.Tick(time.Minute)
	//os.Exit(0)
	//fmt.Println("Delete the rs")
	//msg, _ = json.Marshal(entity.ReplicaSetUpdate{
	//	Action: entity.DeleteAction,
	//	Target: rs,
	//})
	//listwatch.Publish(topic, msg)
	//}()

	cm := controller.NewControllerManager()
	cm.Start()
}
