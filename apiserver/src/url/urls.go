package url

const (
	HttpScheme    = "http://"
	Hostname      = "localhost"
	Port          = ":8080"
	Prefix        = HttpScheme + Hostname + Port
	PodURL        = "/api/v1/pods/"
	ReplicaSetURL = "/api/v1/replicaSet/"
	HPAURL        = "/api/v1/hpa/"
)
