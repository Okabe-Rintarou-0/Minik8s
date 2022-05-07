package util

const podUpdateTopic = "PodUpdate"
const podStatusTopic = "podStatus"
const replicaSetUpdateTopic = "ReplicaSet"
const testTopic = "__test__"

func PodUpdateTopic(hostname string) string {
	return podUpdateTopic + "-" + hostname
}

func SchedulerPodUpdateTopic() string {
	return podUpdateTopic
}

func ReplicaSetUpdateTopic() string {
	return replicaSetUpdateTopic
}

func PodStatusTopic() string {
	return podStatusTopic
}

func TestTopic() string {
	return testTopic
}
