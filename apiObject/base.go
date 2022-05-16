package apiObject

type Metadata struct {
	Name        string      `yaml:"name"`
	Namespace   string      `yaml:"namespace"`
	UID         string      `yaml:"uid"`
	Labels      Labels      `yaml:"labels"`
	Annotations Annotations `yaml:"annotations"`
}
type Base struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
}
