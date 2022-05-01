package podutil

import (
	"minik8s/kubelet/src/types"
	"strconv"
	"strings"
)

func PodFullName(name, namespace string) string {
	return name + "_" + namespace
}

func ContainerFullName(containerName, podFullName string, podUID string, restartCount int) string {
	return "k8s_" + containerName + "_" + podFullName + "_" + podUID + "_" + strconv.Itoa(restartCount)
}

func ToContainerReference(containerFullName string) string {
	return "container:" + containerFullName
}

func ParseContainerFullName(containerFullName string) (succ bool, containerName, podName, podNamespace string, podUID types.UID, restartCount int) {
	if containerFullName[0] == '/' {
		containerFullName = containerFullName[1:]
	}
	tokens := strings.Split(containerFullName, "_")
	var err error
	succ = false
	if numTokens := len(tokens); numTokens == 6 {
		succ = true
		containerName, podName, podNamespace, podUID = tokens[1], tokens[2], tokens[3], tokens[4]
		restartCount, err = strconv.Atoi(tokens[5])
		succ = err == nil
	}
	return
}
