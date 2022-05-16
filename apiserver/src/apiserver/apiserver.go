package apiserver

import (
	"github.com/gin-gonic/gin"
	"log"
	"minik8s/apiserver/src/etcd"
	"minik8s/listwatch"
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

func (api *apiServer) Run() {
	go listwatch.Watch(topicutil.NodeStatusTopic(), syncNodeStatus)
	go listwatch.Watch(topicutil.PodStatusTopic(), syncPodStatus)

	etcd.Start()
	api.bindHandlers()
	log.Fatal(api.httpServer.Run())
}
