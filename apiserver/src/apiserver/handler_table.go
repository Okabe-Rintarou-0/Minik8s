package apiserver

import (
	"github.com/gin-gonic/gin"
	"minik8s/apiserver/src/handlers"
	"minik8s/apiserver/src/url"
)

type Handler = gin.HandlerFunc

var postTable = map[string]Handler{
	// kubectl apply -f xxx.yaml
	url.PodURL:        handlers.HandleApplyPod,
	url.ReplicaSetURL: handlers.HandleApplyReplicaSet,
	url.HPAURL:        handlers.HandleApplyHPA,

	// kubectl autoscale hpa_name -t target -c cpu -m memory --min=min_replicas --max=max_replicas
	url.AutoscaleURLWithSpecifiedName: handlers.HandleAutoscale,

	// kubectl label nodes node_name os=linux --overwrite
	url.NodeLabelsURLWithSpecifiedName: handlers.HandleLabelNode,

	// Set Node Status
	url.NodeStatusURLWithSpecifiedName: handlers.HandleSetNodeStatus,
}

var getTable = map[string]Handler{
	// kubectl get nodes & kubectl get node hostname
	url.NodeURL:                  handlers.HandleGetNodes,
	url.NodeURLWithSpecifiedName: handlers.HandleGetNode,

	// kubectl get pods & kubectl get pod pod_name
	url.PodURL:                  handlers.HandleGetPods,
	url.PodURLWithSpecifiedName: handlers.HandleGetPod,

	// kubectl describe pod pod_name
	url.PodDescriptionURLWithSpecifiedName: handlers.HandleDescribePod,
}

var putTable = map[string]Handler{}

var deleteTable = map[string]Handler{
	// kubectl delete pod pod_name
	url.PodURLWithSpecifiedName: handlers.HandleDeletePod,
}
