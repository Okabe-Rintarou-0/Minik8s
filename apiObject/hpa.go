package apiObject

import "minik8s/apiObject/types"

type ScaleTargetRef = ApiObjectBase

type Metrics struct {
	CPUUtilizationPercentage float64 `yaml:"CPUUtilizationPercentage"`
	MemUtilizationPercentage float64 `yaml:"MemUtilizationPercentage"`
}

type HPASpec struct {
	MinReplicas    int            `yaml:"minReplicas"`
	MaxReplicas    int            `yaml:"maxReplicas"`
	ScaleTargetRef ScaleTargetRef `yaml:"scaleTargetRef"`
	Metrics        *Metrics       `yaml:"metrics"`
}

type HorizontalPodAutoscaler struct {
	ApiObjectBase `yaml:",inline"`
	Spec          HPASpec `yaml:"spec"`
}

func (hpa *HorizontalPodAutoscaler) Name() string {
	return hpa.Metadata.Name
}

func (hpa *HorizontalPodAutoscaler) UID() types.UID {
	return hpa.Metadata.UID
}

func (hpa *HorizontalPodAutoscaler) Namespace() string {
	return hpa.Metadata.Namespace
}

func (hpa *HorizontalPodAutoscaler) Labels() Labels {
	return hpa.Metadata.Labels
}

func (hpa *HorizontalPodAutoscaler) Metrics() *Metrics {
	return hpa.Spec.Metrics
}

func (hpa *HorizontalPodAutoscaler) TargetMetadata() *Metadata {
	return &hpa.Spec.ScaleTargetRef.Metadata
}
