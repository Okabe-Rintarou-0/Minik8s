package url

const (
	HttpScheme = "http://"
	Hostname   = "localhost"
	Port       = ":8080"
	Prefix     = HttpScheme + Hostname + Port

	PodURL                             = "/api/v1/pods/"
	PodURLWithSpecifiedNode            = "/api/v1/pods/nodes/:node"
	PodDescriptionURL                  = "/api/v1/pods/description/"
	PodDescriptionURLWithSpecifiedName = "/api/v1/pods/description/:name"
	PodStatusURLWithSpecifiedName      = "/api/v1/pods/status/:namespace/:name"
	PodURLWithSpecifiedNodeAndName     = "/api/v1/pods/nodes/:node/:namespace/:name"
	PodURLWithSpecifiedName            = "/api/v1/pods/:namespace/:name"

	NodeURL                        = "/api/v1/nodes/"
	NodeURLWithSpecifiedName       = "/api/v1/nodes/:namespace/:name"
	NodeStatusURLWithSpecifiedName = "/api/v1/nodes/status/:namespace/:name/"
	NodeLabelsURLWithSpecifiedName = "/api/v1/nodes/:namespace/:name/labels"

	ReplicaSetURL                        = "/api/v1/replicaSets/"
	ReplicaSetURLWithSpecifiedName       = "/api/v1/replicaSets/:namespace/:name"
	ReplicaSetStatusURLWithSpecifiedName = "/api/v1/replicaSets/status/:namespace/:name"

	HPAURL                        = "/api/v1/hpa/"
	HPAURLWithSpecifiedName       = "/api/v1/hpa/:namespace/:name"
	HPAStatusURLWithSpecifiedName = "/api/v1/hpa/status/:namespace/:name"
	AutoscaleURL                  = "/autoscaling/v1/"
	AutoscaleURLWithSpecifiedName = "/autoscaling/v1/:namespace/:name"

	GpuURL = "/api/v1/gpu/"

	ResetURL = "/reset"
)
