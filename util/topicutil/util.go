package topicutil

const podUpdateTopic = "PodUpdate"
const hpaUpdateTopic = "HPAUpdate"
const podStatusTopic = "PodStatus"
const replicaSetStatusTopic = "ReplicaSetStatus"
const replicaSetUpdateTopic = "ReplicaSet"
const testTopic = "__test__"

func PodUpdateTopic(hostname string) string {
	return podUpdateTopic + "-" + hostname
}

func PodStatusTopic() string {
	return podStatusTopic
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

func TestTopic() string {
	return testTopic
}
