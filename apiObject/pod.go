package apiObject

const (
	PullPolicyAlways = "Always"
)

type Labels map[string]string
type Annotations map[string]string

func (labels Labels) DeepCopy() Labels {
	cpy := make(Labels)
	for key, value := range labels {
		cpy[key] = value
	}
	return cpy
}

type ContainerPort struct {
	Name          string `yaml:"name"`
	HostPort      string `yaml:"hostPort"`
	ContainerPort string `yaml:"containerPort"`
	Protocol      string `yaml:"protocol"`
	HostIP        string `yaml:"hostIP"`
}

type ProbeHandler struct{}

// Probe describes a health check to be performed against a container to determine whether it is
// alive or ready to receive traffic.
type Probe struct {
	// The action taken to determine the health of a container
	ProbeHandler `yaml:",inline"`
	// Number of seconds after the container has started before liveness probes are initiated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	InitialDelaySeconds int32 `yaml:"initial_delay_seconds"`
	// Number of seconds after which the probe times out.
	// Defaults to 1 second. Minimum value is 1.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	TimeoutSeconds int32 `yaml:"timeout_seconds"`
	// How often (in seconds) to perform the probe.
	// Default to 10 seconds. Minimum value is 1.
	// +optional
	PeriodSeconds int32 `yaml:"period_seconds"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
	// +optional
	SuccessThreshold int32 `yaml:"success_threshold"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// Defaults to 3. Minimum value is 1.
	// +optional
	FailureThreshold int32 `yaml:"failure_threshold"`
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
		Command []string `yaml:"cmd,flow"`
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
	Command         []string                     `yaml:"cmd,flow"`
	Args            []string                     `yaml:"args,flow"`
	Env             []EnvVar                     `yaml:"env"`
	Resources       ContainerResources           `yaml:"resources"`
	Ports           []ContainerPort              `yaml:"ports"`
	LivenessProbe   ContainerLivenessProbeConfig `yaml:"livenessProbe"`
	Lifecycle       ContainerLifecycleConfig     `yaml:"lifecycle"`
	VolumeMounts    []VolumeMount                `yaml:"volumeMounts"`
	TTY             bool                         `yaml:"tty"`
}

type EmptyDirVolumeSource struct{}

type HostPathVolumeSource struct {
	Path string `yaml:"path"`
}
type VolumeSource struct {
	EmptyDir *EmptyDirVolumeSource `yaml:"emptyDir"`
	HostPath *HostPathVolumeSource `yaml:"hostPath"`
}

func (vs *VolumeSource) IsEmptyDir() bool {
	return vs.EmptyDir != nil
}

func (vs *VolumeSource) IsHostPath() bool {
	return vs.HostPath != nil
}

type Volume struct {
	Name         string `yaml:"name"`
	VolumeSource `yaml:",inline"`
}

type PodSpec struct {
	RestartPolicy string      `yaml:"restartPolicy"`
	NodeSelector  Labels      `yaml:"nodeSelector,omitempty"`
	Containers    []Container `yaml:"containers"`
	Volumes       []Volume    `yaml:"volumes"`
}

type Pod struct {
	Base `yaml:",inline"` // 处理继承关系，详情请见: https://github.com/go-yaml/yaml/pull/94/commits/e90bcf783f7abddaa0ee0994a09e536498744e49
	Spec PodSpec          `yaml:"spec"`
}

func (pod *Pod) FullName() string {
	return pod.Metadata.Name + "_" + pod.Metadata.Namespace
}

func (pod *Pod) UID() string {
	return pod.Metadata.UID
}

func (pod *Pod) Name() string {
	return pod.Metadata.Name
}

func (pod *Pod) Namespace() string {
	return pod.Metadata.Namespace
}

func (pod *Pod) Labels() Labels {
	return pod.Metadata.Labels
}

func (pod *Pod) Containers() []Container {
	return pod.Spec.Containers
}

func (pod *Pod) AddLabel(name, value string) {
	if pod.Labels() == nil {
		pod.Metadata.Labels = make(Labels)
	}
	pod.Metadata.Labels[name] = value
}

func (pod *Pod) GetContainerByName(name string) *Container {
	for _, container := range pod.Spec.Containers {
		if container.Name == name {
			return &container
		}
	}
	return nil
}

func (pod *Pod) NodeSelector() map[string]string {
	return pod.Spec.NodeSelector
}

type PodTemplateSpec struct {
	Metadata Metadata `yaml:"metadata"`
	Spec     PodSpec  `yaml:"spec"`
}
