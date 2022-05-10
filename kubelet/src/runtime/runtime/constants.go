package runtime

const (
	pauseImage                   = "registry.aliyuncs.com/google_containers/pause:3.6"
	pauseContainerName           = "POD"
	KubernetesPodNameLabel       = "io.kubernetes.pod.name"
	KubernetesPodNamespaceLabel  = "io.kubernetes.pod.namespace"
	KubernetesPodUIDLabel        = "io.kubernetes.pod.uid"
	KubernetesReplicaSetUIDLabel = "io.kubernetes.rs.uid"
	KubernetesContainerNameLabel = "io.kubernetes.container.name"
)
