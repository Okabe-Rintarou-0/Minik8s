package apiObject

type ServicePort struct {
	Name       string `yaml:"name"`
	Port       int32  `yaml:"port,omitempty"`
	TargetPort int32  `yaml:"targetPort,omitempty"`
}

type ServiceSpec struct {
	Type      string            `yaml:"type,omitempty"`
	Ports     []ServicePort     `yaml:"ports,omitempty"`
	Selector  map[string]string `yaml:"selector,omitempty"`
	ClusterIP string            `yaml:"clusterIP,omitempty"`
}

type ServiceStatus struct {
}

type Service struct {
	Base   `yaml:",inline"`
	Spec   ServiceSpec   `yaml:"spec"`
	Status ServiceStatus `yaml:"status,omitempty"`
}

func (service *Service) FullName() string {
	return service.Metadata.Name + "_" + service.Metadata.Namespace
}
