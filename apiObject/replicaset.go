package apiObject

import "minik8s/apiObject/types"

type LabelSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}

// ReplicaSetSpec is the specification of a ReplicaSet.
// As the internal representation of a ReplicaSet, it must have
// a Template set.
type ReplicaSetSpec struct {
	Replicas int             `yaml:"replicas"`
	Selector LabelSelector   `yaml:"selector"`
	Template PodTemplateSpec `yaml:"template"`
}

// ReplicaSet ensures that a specified number of pod replicas are running at any given time.
type ReplicaSet struct {
	Base `yaml:",inline"`
	// Spec defines the desired behavior of this ReplicaSet.
	// +optional
	Spec ReplicaSetSpec `yaml:"spec"`
}

func (template *PodTemplateSpec) ToPod() *Pod {
	return &Pod{
		Base: Base{
			ApiVersion: "v1",
			Kind:       "Pod",
			Metadata:   template.Metadata,
		},
		Spec: template.Spec,
	}
}

func (rs *ReplicaSet) Template() PodTemplateSpec {
	return rs.Spec.Template
}

func (rs *ReplicaSet) UID() types.UID {
	return rs.Metadata.UID
}

func (rs *ReplicaSet) Replicas() int {
	return rs.Spec.Replicas
}

func (rs *ReplicaSet) SetReplicas(numReplicas int) {
	rs.Spec.Replicas = numReplicas
}

func (rs *ReplicaSet) Name() string {
	return rs.Metadata.Name
}

func (rs *ReplicaSet) Namespace() string {
	return rs.Metadata.Namespace
}

func (rs *ReplicaSet) Labels() Labels {
	return rs.Metadata.Labels
}

func (rs *ReplicaSet) FullName() string {
	return rs.Metadata.Name + "_" + rs.Metadata.Namespace
}
