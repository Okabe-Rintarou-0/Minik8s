package apiObject

import "minik8s/apiObject/types"

type Node struct {
	Base `yaml:",inline"`
	Ip   string `yaml:"ip,omitempty"`
}

func (node *Node) UID() types.UID {
	return node.Metadata.UID
}

func (node *Node) Name() string {
	return node.Metadata.Name
}

func (node *Node) Namespace() string {
	return node.Metadata.Namespace
}

func (node *Node) Labels() Labels {
	return node.Metadata.Labels
}
