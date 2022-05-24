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

func ipInit(url, ipBase string) error {
	if ig, err := ipgen.New(url, ipBase); err != nil {
		return err
	} else {
		if err := ig.ClearIfInit(); err != nil {
			return err
		}
	}
	return nil
}

func weaveInit(url, ipBase string) error {
	if out, err := exec.Command("weave", "reset").Output(); err != nil {
		return err
	} else {
		logger.Log("api-server")(string(out))
	}
	if out, err := exec.Command("weave", "launch").Output(); err != nil {
		return err
	} else {
		logger.Log("api-server")(string(out))
	}
	if ig, err := ipgen.New(url, ipBase); err != nil {
		return err
	} else {
		if ip, err := ig.GetNextWithMask(); err != nil {
			return err
		} else {
			if out, err := exec.Command("weave", "expose", ip).Output(); err != nil {
				return err
			} else {
				logger.Log("api-server")(string(out))
			}
		}
	}
	return nil
}

func (api *apiServer) Run() {
	etcd.Start()
	_ = etcd.DeleteAllKeys()

	if err := ipInit(url.SvcIpGeneratorURL, url.ServiceIpBase); err != nil {
		logger.Log("api-server")(err.Error())
		return
	}
	if err := ipInit(url.PodIpGeneratorURL, url.PodIpBase); err != nil {
		logger.Log("api-server")(err.Error())
		return
	}
	if err := weaveInit(url.PodIpGeneratorURL, url.PodIpBase); err != nil {
		logger.Log("api-server")(err.Error())
		return
	}
	api.bindHandlers()
	api.watch()
	log.Fatal(api.httpServer.Run())
}
