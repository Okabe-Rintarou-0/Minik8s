package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/kubectl/src/util"
	"minik8s/util/apiutil"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Kubectl apply is used to create api object in a declarative way",
	Long:  "Kubectl apply is used to create api object in a declarative way",
	Run:   apply,
}

func apply(cmd *cobra.Command, args []string) {
	content, err := util.LoadContent(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tp, err := util.ParseApiObjectType(content)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch tp {
	case util.Node:
		node := apiObject.Node{}
		if err = yaml.Unmarshal(content, &node); err != nil {
			fmt.Println(err.Error())
			return
		}
		URL := url.Prefix + url.NodeURL
		apiutil.ApplyApiObjectToApiServer(URL, node)
	case util.Pod:
		pod := apiObject.Pod{}
		if err = yaml.Unmarshal(content, &pod); err != nil {
			fmt.Println(err.Error())
			return
		}
		URL := url.Prefix + url.PodURL
		apiutil.ApplyApiObjectToApiServer(URL, pod)
	case util.ReplicaSet:
		rs := apiObject.ReplicaSet{}
		if err = yaml.Unmarshal(content, &rs); err != nil {
			fmt.Println(err.Error())
			return
		}
		URL := url.Prefix + url.ReplicaSetURL
		apiutil.ApplyApiObjectToApiServer(URL, rs)
	case util.HorizontalPodAutoscaler:
		hpa := apiObject.HorizontalPodAutoscaler{}
		if err = yaml.Unmarshal(content, &hpa); err != nil {
			fmt.Println(err.Error())
			return
		}
		if hpa.MinReplicas() > hpa.MaxReplicas() {
			fmt.Println("Minimum number of replicas should be less than Maximum one!")
			return
		}
		URL := url.Prefix + url.HPAURL
		apiutil.ApplyApiObjectToApiServer(URL, hpa)
	case util.Service:
		service := apiObject.Service{}
		if err = yaml.Unmarshal(content, &service); err != nil {
			fmt.Println(err.Error())
			return
		}
		URL := url.Prefix + url.ServiceURL
		apiutil.ApplyApiObjectToApiServer(URL, service)
	case util.DNS:
		dns := apiObject.Dns{}
		if err = yaml.Unmarshal(content, &dns); err != nil {
			fmt.Println(err.Error())
			return
		}
		URL := url.Prefix + url.DNSURL
		apiutil.ApplyApiObjectToApiServer(URL, dns)
	case util.GpuJob:
		gpu := apiObject.GpuJob{}
		if err = yaml.Unmarshal(content, &gpu); err != nil {
			fmt.Println(err.Error())
			return
		}
		URL := url.Prefix + url.GpuURL
		apiutil.ApplyApiObjectToApiServer(URL, gpu)
	}
}
