package apiObject

import "minik8s/apiObject/types"

type ScaleTargetRef = Base

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
	Base `yaml:",inline"`
	Spec HPASpec `yaml:"spec"`
}

func (tgt *ScaleTargetRef) Name() string {
	return tgt.Metadata.Name
}

func (tgt *ScaleTargetRef) Namespace() string {
	return tgt.Metadata.Namespace
}

func (tgt *ScaleTargetRef) FullName() string {
	return tgt.Metadata.Name + "_" + tgt.Metadata.Namespace
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

func (hpa *HorizontalPodAutoscaler) MinReplicas() int {
	return hpa.Spec.MinReplicas
}

func (hpa *HorizontalPodAutoscaler) MaxReplicas() int {
	return hpa.Spec.MaxReplicas
}

func (hpa *HorizontalPodAutoscaler) Target() *ScaleTargetRef {
	return &hpa.Spec.ScaleTargetRef
}

func (hpa *HorizontalPodAutoscaler) TargetMetadata() *Metadata {
	return &hpa.Spec.ScaleTargetRef.Metadata
}
