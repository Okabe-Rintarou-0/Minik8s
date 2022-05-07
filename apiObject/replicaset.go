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
	ApiObjectBase `yaml:",inline"`
	// Spec defines the desired behavior of this ReplicaSet.
	// +optional
	Spec ReplicaSetSpec `yaml:"spec"`
}

func (rs *ReplicaSet) UID() types.UID {
	return rs.Metadata.UID
}

func (rs *ReplicaSet) Replicas() int {
	return rs.Spec.Replicas
}

func (rs *ReplicaSet) FullName() string {
	return rs.Metadata.Name + "_" + rs.Metadata.Namespace
}
