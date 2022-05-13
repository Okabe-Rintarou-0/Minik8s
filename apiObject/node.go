package apiObject

type Node struct {
	ApiObjectBase `yaml:",inline"`
	Ip            string `yaml:"ip,omitempty"`
}
