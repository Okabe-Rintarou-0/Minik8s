package parseutil

import (
	"gopkg.in/yaml.v3"
	"minik8s/apiObject"
)

func ParsePod(content []byte) (*apiObject.Pod, error) {
	pod := &apiObject.Pod{}
	err := yaml.Unmarshal(content, pod)
	return pod, err
}

func ParseHPA(content []byte) (*apiObject.HorizontalPodAutoscaler, error) {
	hpa := &apiObject.HorizontalPodAutoscaler{}
	err := yaml.Unmarshal(content, hpa)
	return hpa, err
}
