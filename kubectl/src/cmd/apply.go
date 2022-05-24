package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/kubectl/src/util"
	"minik8s/util/httputil"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Kubectl apply is used to create api object in a declarative way",
	Long:  "Kubectl apply is used to create api object in a declarative way",
	Run:   apply,
}

func applyApiObjectToApiServer(URL string, object interface{}) {
	resp, err := httputil.PostJson(URL, object)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(respBody))
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
		applyApiObjectToApiServer(URL, node)
	case util.Pod:
		pod := apiObject.Pod{}
		if err = yaml.Unmarshal(content, &pod); err != nil {
			fmt.Println(err.Error())
			return
		}
		URL := url.Prefix + url.PodURL
		applyApiObjectToApiServer(URL, pod)
	case util.ReplicaSet:
		rs := apiObject.ReplicaSet{}
		if err = yaml.Unmarshal(content, &rs); err != nil {
			fmt.Println(err.Error())
			return
		}
		URL := url.Prefix + url.ReplicaSetURL
		applyApiObjectToApiServer(URL, rs)
	case util.HorizontalPodAutoscaler:
		hpa := apiObject.HorizontalPodAutoscaler{}
		if err = yaml.Unmarshal(content, &hpa); err != nil {
			fmt.Println(err.Error())
			return
		}
		URL := url.Prefix + url.HPAURL
		applyApiObjectToApiServer(URL, hpa)
	case util.Service:
		service := apiObject.Service{}
		if err = yaml.Unmarshal(content, &service); err != nil {
			fmt.Println(err.Error())
			return
		}
		URL := url.Prefix + url.ServiceURL
		applyApiObjectToApiServer(URL, service)
	}
}
