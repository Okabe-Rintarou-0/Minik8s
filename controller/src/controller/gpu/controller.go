package gpu

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/httputil"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"path"
)

var log = logger.Log("Gpu")

const (
	minik8sGpuServerImage = "923048992/minik8s-gpu-server:v1.2"
)

type Controller interface {
	Run()
}

type controller struct{}

func (c *controller) dispatchGpuJob(gpuJob *apiObject.GpuJob) {
	log("received %+v gpu: ", gpuJob)
	jobFullName := gpuJob.Namespace() + "/" + gpuJob.Name()
	gpuServerCommands := []string{
		"./gpu-server",
		fmt.Sprintf("--job-name=%s", jobFullName),
		fmt.Sprintf("--workdir=%s", gpuJob.WorkDir()),
		fmt.Sprintf("--output=%s", gpuJob.OutputFile()),
		fmt.Sprintf("--error=%s", gpuJob.ErrorFile()),
		fmt.Sprintf("--process=%d", gpuJob.NumProcess()),
		fmt.Sprintf("--ntasks-per-node=%d", gpuJob.NumTasksPerNode()),
		fmt.Sprintf("--cpus-per-task=%d", gpuJob.CpusPerTask()),
		fmt.Sprintf("--gres=gpu:%d", gpuJob.NumGpus()),
		fmt.Sprintf("--username=%s", gpuJob.Username()),
		fmt.Sprintf("--password=%s", gpuJob.Password()),
		fmt.Sprintf("--run=%s", gpuJob.RunScripts()),
		fmt.Sprintf("--compile=%s", gpuJob.CompileScripts()),
	}

	podNamePrefix := jobFullName
	pod := apiObject.Pod{
		Base: apiObject.Base{
			ApiVersion: "v1",
			Kind:       "Pod",
			Metadata: apiObject.Metadata{
				Name:      gpuJob.Name(),
				Namespace: gpuJob.Namespace(),
			},
		},
		Spec: apiObject.PodSpec{
			Containers: []apiObject.Container{
				{
					Name:  podNamePrefix + "-file-server",
					Image: "dplsming/nginx-fileserver:1.0",
					Ports: []apiObject.ContainerPort{
						{
							ContainerPort: "80",
						},
					},
					VolumeMounts: []apiObject.VolumeMount{
						{
							Name:      "volume",
							MountPath: "/usr/share/nginx/html/files",
						},
					},
				},
				{
					Name:    podNamePrefix + "-gpu-server",
					Image:   minik8sGpuServerImage,
					Command: gpuServerCommands,
					VolumeMounts: []apiObject.VolumeMount{
						{
							Name:      "volume",
							MountPath: "/usr/local/jobs",
						},
					},
				},
			},
			Volumes: []apiObject.Volume{
				{
					Name: "volume",
					VolumeSource: apiObject.VolumeSource{
						HostPath: &apiObject.HostPathVolumeSource{
							Path: gpuJob.Volume(),
						},
					},
				},
			},
		},
	}

	URL := url.Prefix + url.PodURL
	if resp, err := httputil.PostJson(URL, pod); err == nil {
		content, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		log("Apply pod and got resp: %s", content)
	} else {
		logger.Error(err.Error())
	}
}

func (c *controller) deleteGpuJobPod(gpuJob *apiObject.GpuJob) {
	namespace, name := gpuJob.Namespace(), gpuJob.Name()
	URL := url.Prefix + path.Join(url.PodURL, namespace, name)
	httputil.DeleteWithoutBody(URL)
}

func (c *controller) handleGpuJobUpdate(msg *redis.Message) {
	gpuJobUpdate := &entity.GpuUpdate{}
	if err := json.Unmarshal([]byte(msg.Payload), gpuJobUpdate); err != nil {
		return
	}
	gpuJob := gpuJobUpdate.Target
	switch gpuJobUpdate.Action {
	case entity.CreateAction:
		c.dispatchGpuJob(&gpuJob)
	case entity.DeleteAction:
		c.deleteGpuJobPod(&gpuJob)
	}
}

func (c *controller) Run() {
	listwatch.Watch(topicutil.GpuJobUpdateTopic(), c.handleGpuJobUpdate)
}

func NewController() Controller {
	return &controller{}
}
