package util

const podUpdateTopic = "PodUpdate"

func PodUpdateTopic(hostname string) string {
	return podUpdateTopic + "-" + hostname
}

func SchedulerPodUpdateTopic() string {
	return podUpdateTopic
}
