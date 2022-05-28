package apiserver

import (
	"github.com/gin-gonic/gin"
	"minik8s/apiserver/src/handlers"
	"minik8s/apiserver/src/url"
)

type Handler = gin.HandlerFunc

var postTable = map[string]Handler{
	// kubectl apply -f xxx.yaml
	url.NodeURL:       handlers.HandleApplyNode,
	url.PodURL:        handlers.HandleApplyPod,
	url.ReplicaSetURL: handlers.HandleApplyReplicaSet,
	url.HPAURL:        handlers.HandleApplyHPA,
	url.ServiceURL:    handlers.HandleApplyService,
	url.DNSURL:        handlers.HandleApplyDNS,
	url.GpuURL:        handlers.HandleApplyGpuJob,

	// kubectl wf apply -f
	url.WorkflowURL: handlers.HandleApplyWorkflow,

	// update pod after it's scheduled
	url.PodURLWithSpecifiedNode: handlers.HandleSchedulePod,

	// kubectl autoscale hpa_name -t target -c cpu -m memory --min=min_replicas --max=max_replicas
	url.AutoscaleURLWithSpecifiedName: handlers.HandleAutoscale,

	// kubectl label nodes node_name os=linux --overwrite
	url.NodeLabelsURLWithSpecifiedName: handlers.HandleLabelNode,

	// kubectl func add func_name func_path
	url.FuncURL: handlers.HandleApplyFunc,

	// post workflow result
	url.WorkflowResultURLWithSpecifiedName: handlers.HandlePutWorkflowResult,
}

var getTable = map[string]Handler{
	// kubectl get nodes & kubectl get node hostname
	url.NodeURL:                        handlers.HandleGetNodeStatuses,
	url.NodeStatusURLWithSpecifiedName: handlers.HandleGetNodeStatus,

	// kubectl get pods & kubectl get pod pod_name
	url.PodURL:                        handlers.HandleGetPodStatuses,
	url.PodStatusURLWithSpecifiedName: handlers.HandleGetPodStatus,

	// kubectl get replicaSets && kubectl get replicaSet replicaSet_name
	url.ReplicaSetURL:                        handlers.HandleGetReplicaSetStatuses,
	url.ReplicaSetStatusURLWithSpecifiedName: handlers.HandleGetReplicaSetStatus,

	// kubectl get hpa && kubectl get hpa hpa_name
	url.HPAURL:                        handlers.HandleGetHPAStatuses,
	url.HPAStatusURLWithSpecifiedName: handlers.HandleGetHPAStatus,

	// kubectl describe pod pod_name
	url.PodDescriptionURLWithSpecifiedName: handlers.HandleDescribePod,

	// get apiObject.xxx
	url.PodURLWithSpecifiedNodeAndName: handlers.HandleGetPodApiObject,
	url.ReplicaSetURLWithSpecifiedName: handlers.HandleGetReplicaSetApiObject,
	url.HPAURLWithSpecifiedName:        handlers.HandleGetHPAApiObject,
	url.PodURLWithSpecifiedNode:        handlers.HandleGetPodsApiObject,
	url.GpuURLWithSpecifiedName:        handlers.HandleGetGpuApiObject,

	// kubectl get service service_name
	url.ServiceURLWithSpecifiedName: handlers.HandleGetService,
	// kubectl get services
	url.ServiceURL: handlers.HandleGetServices,

	// kubectl get dns dns_name
	url.DNSURLWithSpecifiedName: handlers.HandleGetDNS,
	url.DNSURL:                  handlers.HandleGetDNSes,

	// get pod function pods
	url.FuncPodsURLWithSpecifiedName: handlers.HandleGetFuncPods,

	// get workflow result
	url.WorkflowResultURLWithSpecifiedName: handlers.HandleGetWorkflowResult,

	// get all workflow results
	url.WorkflowURL: handlers.HandleGetWorkflowResults,
}

var putTable = map[string]Handler{
	// Set Node Status
	url.NodeStatusURLWithSpecifiedName: handlers.HandleSetNodeStatus,
	url.ReplicaSetURLWithSpecifiedName: handlers.HandleSetReplicaSet,
}

var deleteTable = map[string]Handler{
	// kubectl delete apiObjectType apiObjectName
	url.PodURLWithSpecifiedName:        handlers.HandleDeletePod,
	url.NodeURLWithSpecifiedName:       handlers.HandleDeleteNode,
	url.ReplicaSetURLWithSpecifiedName: handlers.HandleDeleteReplicaSet,
	url.HPAURLWithSpecifiedName:        handlers.HandleDeleteHPA,
	url.ServiceURLWithSpecifiedName:    handlers.HandleDeleteService,
	url.DNSURLWithSpecifiedName:        handlers.HandleDeleteDNS,
	url.GpuURLWithSpecifiedName:        handlers.HandleDeleteGpuJob,

	// kubectl reset
	url.ResetURL: handlers.HandleReset,

	// delete pods of a node
	url.PodURLWithSpecifiedNode: handlers.HandleDeleteNodePods,

	// kubectl func rm func_name
	url.FuncURLWithSpecifiedName: handlers.HandleRemoveFunc,

	// kubectl wf rm workflow name:
	url.WorkflowURLWithSpecifiedName: handlers.HandleRemoveWorkflow,
}
