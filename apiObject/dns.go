package apiObject

type DnsPath struct {
	Path        string `yaml:"path,omitempty"`
	ServiceName string `yaml:"serviceName,omitempty"`
	ServicePort string `yaml:"servicePort,omitempty"`
}

type DnsSpec struct {
	Host  string    `yaml:"host,omitempty"`
	Paths []DnsPath `yaml:"paths,omitempty"`
}

type DnsStatus struct {
}

type Dns struct {
	Base   `yaml:",inline"`
	Spec   DnsSpec   `yaml:"spec"`
	Status DnsStatus `yaml:"status,omitempty"`
}
