package apiObject

type Labels map[string]string
type Annotations map[string]string

type ContainerPortConfig struct {
	ContainerPort int    `yaml:"containerPort"`
	Name          string `yaml:"name"`
	Protocol      string `yaml:"protocol"`
}

type ContainerEnvironmentVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type ContainerResources struct {
	Cpu    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type ContainerResourcesConfig struct {
	Requests ContainerResources `yaml:"requests"`
	Limits   ContainerResources `yaml:"limits"`
}

type ContainerLivenessProbeConfig struct {
	HttpGet struct {
		Path   string `yaml:"path"`
		Port   int    `yaml:"port"`
		Host   string `yaml:"host"`
		Scheme string `yaml:"scheme"`
	} `yaml:"httpGet"`
	InitialDelaySeconds int `yaml:"initialDelaySeconds"`
	TimeoutSeconds      int `yaml:"timeoutSeconds"`
	PeriodSeconds       int `yaml:"periodSeconds"`
}

type ContainerLifecycleTaskConfig struct {
	Exec struct {
		Command []string `yaml:"command,flow"`
	} `yaml:"exec"`
}

type ContainerLifecycleConfig struct {
	PostStart ContainerLifecycleTaskConfig `yaml:"postStart"`
	PreStop   ContainerLifecycleTaskConfig `yaml:"preStop"`
}

type ContainerVolumeMountConfig struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
	ReadOnly  bool   `yaml:"readOnly"`
}

type ContainerConfig struct {
	Name            string                         `yaml:"name"`
	Image           string                         `yaml:"image"`
	ImagePullPolicy string                         `yaml:"imagePullPolicy"`
	Command         []string                       `yaml:"command,flow"`
	Args            []string                       `yaml:"args,flow"`
	Env             []ContainerEnvironmentVariable `yaml:"env"`
	Resources       ContainerResourcesConfig       `yaml:"resources"`
	Ports           []ContainerPortConfig          `yaml:"ports"`
	LivenessProbe   ContainerLivenessProbeConfig   `yaml:"livenessProbe"`
	Lifecycle       ContainerLifecycleConfig       `yaml:"lifecycle"`
	VolumeMounts    []ContainerVolumeMountConfig   `yaml:"volumeMounts"`
}

type PodVolumeEmptyDirConfig struct {
}

type PodVolumeHostPathConfig struct {
	Path string `yaml:"path"`
}

type PodVolumeConfig struct {
	Name     string                  `yaml:"name"`
	EmptyDir PodVolumeEmptyDirConfig `yaml:"emptyDir"`
	HostPath PodVolumeHostPathConfig `yaml:"hostPath"`
}

type NodeSelector struct {
	Zone    string `yaml:"zone"`
	IpRange string `yaml:"ipRange"`
}

type PodSpec struct {
	RestartPolicy string            `yaml:"restartPolicy"`
	NodeSelector  NodeSelector      `yaml:"nodeSelector"`
	Containers    []ContainerConfig `yaml:"containers"`
	Volumes       []PodVolumeConfig `yaml:"volumes"`
}

type Pod struct {
	ApiObjectBase `yaml:",inline"` // 处理继承关系，详情请见: https://github.com/go-yaml/yaml/pull/94/commits/e90bcf783f7abddaa0ee0994a09e536498744e49
	Spec          PodSpec          `yaml:"spec"`
}
