package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/kubectl/src/util"
	"minik8s/util/httputil"
	"minik8s/util/parseutil"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Kubectl apply is used to create api object in a declarative way",
	Long:  "Kubectl apply is used to create api object in a declarative way",
	Run:   apply,
}

func applyPodToApiServer(pod *apiObject.Pod) {
	resp, err := httputil.PostJson(url.Prefix+url.PodURL, pod)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(respBody))
}

func applyReplicaSetToApiServer(rs *apiObject.ReplicaSet) {
	resp, err := httputil.PostJson(url.Prefix+url.ReplicaSetURL, rs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(respBody))
}

func applyHPAToApiServer(hpa *apiObject.HorizontalPodAutoscaler) {
	resp, err := httputil.PostJson(url.Prefix+url.HPAURL, hpa)
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
	case util.Pod:
		pod, err := parseutil.ParsePod(content)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		applyPodToApiServer(pod)
	case util.ReplicaSet:
		rs, err := parseutil.ParseReplicaSet(content)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		applyReplicaSetToApiServer(rs)
	case util.HorizontalPodAutoscaler:
		hpa, err := parseutil.ParseHPA(content)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		applyHPAToApiServer(hpa)
	}
}
