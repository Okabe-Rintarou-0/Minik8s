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
const gpuJobUpdateTopic = "GpuJobUpdate"
const functionUpdateTopic = "FunctionUpdate"
const scheduleStrategyTopic = "ScheduleStrategyTopic"
const functionTriggerTopic = "FunctionTriggerTopic"
const workflowUpdateTopic = "WorkflowUpdateTopic"
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

func GpuJobUpdateTopic() string {
	return gpuJobUpdateTopic
}

func WorkflowUpdateTopic() string {
	return workflowUpdateTopic
}

func FunctionUpdateTopic() string {
	return functionUpdateTopic
}

func FunctionTriggerTopic() string {
	return functionTriggerTopic
}

func ScheduleStrategyTopic() string {
	return scheduleStrategyTopic
}

func TestTopic() string {
	return testTopic
}
