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
	for URL, handler := range postTable {
		api.httpServer.POST(URL, handler)
	}

	for URL, handler := range getTable {
		api.httpServer.GET(URL, handler)
	}

	for URL, handler := range deleteTable {
		api.httpServer.DELETE(URL, handler)
	}

	for URL, handler := range putTable {
		api.httpServer.PUT(URL, handler)
	}
}

func (api *apiServer) watch() {
	go listwatch.Watch(topicutil.NodeStatusTopic(), syncNodeStatus)
	go listwatch.Watch(topicutil.PodStatusTopic(), syncPodStatus)
	go listwatch.Watch(topicutil.ReplicaSetStatusTopic(), syncReplicaSetStatus)
	go listwatch.Watch(topicutil.HPAStatusTopic(), syncHPAStatus)
	go listwatch.Watch(topicutil.GpuJobUpdateTopic(), syncGpuJobStatus)
}

func ipInit(url, ip string, mask int) error {
	ig := ipgen.New(url, mask)
	return ig.ClearIfNotInit(ip)
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
