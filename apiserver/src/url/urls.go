package url

const (
	HttpScheme = "http://"
	Hostname   = "localhost"
	Port       = ":8080"
	Prefix     = HttpScheme + Hostname + Port

	PodURL                             = "/api/v1/pods/"
	PodDescriptionURL                  = "/api/v1/pods/description/"
	PodDescriptionURLWithSpecifiedName = "/api/v1/pods/description/:name"
	PodURLWithSpecifiedName            = "/api/v1/pods/:name"

	NodeURL                        = "/api/v1/nodes/"
	NodeURLWithSpecifiedName       = "/api/v1/nodes/:name"
	NodeStatusURLWithSpecifiedName = "/api/v1/nodes/:name/status"
	NodeLabelsURLWithSpecifiedName = "/api/v1/nodes/:name/labels"

	ReplicaSetURL = "/api/v1/replicaSet/"

	HPAURL                        = "/api/v1/hpa/"
	AutoscaleURL                  = "/autoscaling/v1/"
	AutoscaleURLWithSpecifiedName = "/autoscaling/v1/:name"
)
