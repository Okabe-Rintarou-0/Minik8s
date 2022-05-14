package apiObject

type Node struct {
	Base `yaml:",inline"`
	Ip   string `yaml:"ip,omitempty"`
}
