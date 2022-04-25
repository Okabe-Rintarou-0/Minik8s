package apiObject

const (
	PullPolicyAlways = "Always"
)

type Labels map[string]string
type Annotations map[string]string

type ContainerPort struct {
	Name          string `yaml:"name"`
	HostPort      string `yaml:"hostPort"`
	ContainerPort string `yaml:"containerPort"`
	Protocol      string `yaml:"protocol"`
	HostIP        string `yaml:"hostIP"`
}

type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type ContainerMetrics struct {
	Cpu    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type ContainerResources struct {
	Requests ContainerMetrics `yaml:"requests"`
	Limits   ContainerMetrics `yaml:"limits"`
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

type VolumeMount struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
	ReadOnly  bool   `yaml:"readOnly"`
}

type Container struct {
	Name            string                       `yaml:"name"`
	Image           string                       `yaml:"image"`
	ImagePullPolicy string                       `yaml:"imagePullPolicy"`
	Command         []string                     `yaml:"command,flow"`
	Args            []string                     `yaml:"args,flow"`
	Env             []EnvVar                     `yaml:"env"`
	Resources       ContainerResources           `yaml:"resources"`
	Ports           []ContainerPort              `yaml:"ports"`
	LivenessProbe   ContainerLivenessProbeConfig `yaml:"livenessProbe"`
	Lifecycle       ContainerLifecycleConfig     `yaml:"lifecycle"`
	VolumeMounts    []VolumeMount                `yaml:"volumeMounts"`
	TTY             bool                         `yaml:"tty"`
}

type EmptyDirVolumeSource struct {
}

type HostPathVolumeSource struct {
	Path string `yaml:"path"`
}
type VolumeSource struct {
	EmptyDir *EmptyDirVolumeSource `yaml:"emptyDir"`
	HostPath *HostPathVolumeSource `yaml:"hostPath"`
}
type Volume struct {
	Name         string `yaml:"name"`
	VolumeSource `yaml:",inline"`
}

type PodSpec struct {
	RestartPolicy string            `yaml:"restartPolicy"`
	NodeSelector  map[string]string `yaml:"nodeSelector,omitempty"`
	Containers    []Container       `yaml:"containers"`
	Volumes       []Volume          `yaml:"volumes"`
}

type Pod struct {
	ApiObjectBase `yaml:",inline"` // 处理继承关系，详情请见: https://github.com/go-yaml/yaml/pull/94/commits/e90bcf783f7abddaa0ee0994a09e536498744e49
	Spec          PodSpec          `yaml:"spec"`
}

func (pod *Pod) FullName() string {
	return pod.Metadata.Name + "_" + pod.Metadata.Namespace
}

func (pod *Pod) UID() string {
	return pod.Metadata.UID
}
