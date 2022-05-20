package apiObject

import "minik8s/apiObject/types"

type GpuJobSpec struct {
	NumProcess      int      `yaml:"numProcess"`
	NumTasksPerNode int      `yaml:"numTasksPerNode"`
	CpusPerTask     int      `yaml:"cpusPerTask"`
	NumGpus         int      `yaml:"numGpus"`
	Scripts         []string `yaml:"scripts"`
	Volume          string   `yaml:"volume"`
	OutputFile      string   `yaml:"outputFile"`
	ErrorFile       string   `yaml:"errorFile"`
}

type GpuJob struct {
	Base `yaml:",inline"`
	Spec GpuJobSpec `yaml:"spec"`
}

func (gpu *GpuJob) Namespace() string {
	return gpu.Metadata.Namespace
}

func (gpu *GpuJob) Name() string {
	return gpu.Metadata.Name
}

func (gpu *GpuJob) UID() types.UID {
	return gpu.Metadata.UID
}

func (gpu *GpuJob) Volume() string {
	return gpu.Spec.Volume
}

func (gpu *GpuJob) OutputFile() string {
	return gpu.Spec.OutputFile
}

func (gpu *GpuJob) ErrorFile() string {
	return gpu.Spec.ErrorFile
}

func (gpu *GpuJob) NumProcess() int {
	return gpu.Spec.NumProcess
}

func (gpu *GpuJob) NumTasksPerNode() int {
	return gpu.Spec.NumTasksPerNode
}

func (gpu *GpuJob) CpusPerTask() int {
	return gpu.Spec.CpusPerTask
}

func (gpu *GpuJob) NumGpus() int {
	return gpu.Spec.NumGpus
}
