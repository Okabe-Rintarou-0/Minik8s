package apiserver

import (
	"github.com/gin-gonic/gin"
	"minik8s/apiserver/src/handlers"
	"minik8s/apiserver/src/url"
)

type Handler = gin.HandlerFunc

var postTable = map[string]Handler{
	url.PodURL:        handlers.HandleApplyPod,
	url.ReplicaSetURL: handlers.HandleApplyReplicaSet,
	url.HPAURL:        handlers.HandleApplyHPA,
}

var getTable = map[string]Handler{}

var putTable = map[string]Handler{}

var deleteTable = map[string]Handler{}
