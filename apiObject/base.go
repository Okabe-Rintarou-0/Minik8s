package apiObject

type MetadataConfigBase struct {
	Name        string      `yaml:"name"`
	Namespace   string      `yaml:"namespace"`
	UID         string      `yaml:"uid"`
	Labels      Labels      `yaml:"labels"`
	Annotations Annotations `yaml:"annotations"`
}
type ApiObjectBase struct {
	ApiVersion string             `yaml:"apiVersion"`
	Kind       string             `yaml:"kind"`
	Metadata   MetadataConfigBase `yaml:"metadata"`
}
