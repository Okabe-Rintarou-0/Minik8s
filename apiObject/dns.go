package apiObject

type DnsService struct {
	Name string `yaml:"name,omitempty"`
	Port string `yaml:"port,omitempty"`
}

type DnsPath struct {
	Path    string     `yaml:"path,omitempty"`
	Service DnsService `yaml:"service,omitempty"`
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
