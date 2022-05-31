package apiserver

import (
	"github.com/gin-gonic/gin"
	"log"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/ipgen"
	"minik8s/apiserver/src/url"
	"minik8s/listwatch"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"os/exec"
)

type ApiServer interface {
	Run()
}

func New() ApiServer {
	return &apiServer{
		httpServer: gin.Default(),
	}
}

type apiServer struct {
	httpServer *gin.Engine
}

func (api *apiServer) bindHandlers() {
	for url, handler := range postTable {
		api.httpServer.POST(url, handler)
	}

	for url, handler := range getTable {
		api.httpServer.GET(url, handler)
	}

	for url, handler := range deleteTable {
		api.httpServer.DELETE(url, handler)
	}

	for url, handler := range putTable {
		api.httpServer.PUT(url, handler)
	}
}

func (api *apiServer) watch() {
	go listwatch.Watch(topicutil.NodeStatusTopic(), syncNodeStatus)
	go listwatch.Watch(topicutil.PodStatusTopic(), syncPodStatus)
	go listwatch.Watch(topicutil.ReplicaSetStatusTopic(), syncReplicaSetStatus)
	go listwatch.Watch(topicutil.HPAStatusTopic(), syncHPAStatus)
}

func ipInit(url, ip string, mask int) error {
	ig := ipgen.New(url, mask)
	return ig.ClearIfInit(ip)
}

func weaveInit() error {
	if out, err := exec.Command("weave", "reset").Output(); err != nil {
		return err
	} else {
		logger.Log("api-server-weave")(string(out))
	}
	if out, err := exec.Command("weave", "launch").Output(); err != nil {
		return err
	} else {
		logger.Log("api-server-weave")(string(out))
	}
	if out, err := exec.Command("weave", "expose", url.PodIpBase+url.MaskStr).Output(); err != nil {
		return err
	} else {
		logger.Log("api-server-weave")(string(out))
	}
	return nil
}

func (api *apiServer) Run() {
	etcd.Start()
	//_ = etcd.DeleteAllKeys()

	if err := ipInit(url.PodIpURL, url.PodIpBase, url.Mask); err != nil {
		logger.Log("api-server-pod-ip")(err.Error())
		return
	}
	if err := ipInit(url.ServiceIpURL, url.ServiceIpBase, url.Mask); err != nil {
		logger.Log("api-server-service-ip")(err.Error())
		return
	}
	api.bindHandlers()
	api.watch()
	log.Fatal(api.httpServer.Run(":8080"))
}
