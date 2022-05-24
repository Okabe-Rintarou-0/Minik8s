package topicutil

const podUpdateTopic = "PodUpdate"
const hpaUpdateTopic = "HPAUpdate"
const hpaStatusTopic = "HPAStatus"
const podStatusTopic = "PodStatus"
const nodeStatusTopic = "NodeStatus"
const replicaSetStatusTopic = "ReplicaSetStatus"
const replicaSetUpdateTopic = "ReplicaSetUpdate"
const serviceUpdateTopic = "ServiceUpdate"
const endpointUpdateTopic = "EndpointUpdate"
const testTopic = "__test__"

func PodUpdateTopic(hostname string) string {
	return podUpdateTopic + "-" + hostname
}

func PodStatusTopic() string {
	return podStatusTopic
}

func NodeStatusTopic() string {
	return nodeStatusTopic
}

func SchedulerPodUpdateTopic() string {
	return podUpdateTopic
}

func ReplicaSetUpdateTopic() string {
	return replicaSetUpdateTopic
}

func ReplicaSetStatusTopic() string {
	return replicaSetStatusTopic
}

func HPAUpdateTopic() string {
	return hpaUpdateTopic
}

func HPAStatusTopic() string {
	return hpaStatusTopic
}

func ServiceUpdateTopic() string {
	return serviceUpdateTopic
}

func EndpointUpdateTopic() string {
	return endpointUpdateTopic
}

func TestTopic() string {
	return testTopic
}
