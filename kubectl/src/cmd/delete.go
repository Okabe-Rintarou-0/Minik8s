package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/apiserver/src/url"
	"minik8s/util/httputil"
	"path"
	"strings"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Kubectl delete is used to delete api object given its name",
	Long:  "Kubectl delete is used to delete api object given its name",
	Args:  cobra.ExactValidArgs(2),
	Run:   del,
}

func deleteSpecifiedNode(namespace, name string) error {
	resp := httputil.DeleteWithoutBody(url.Prefix + path.Join(url.NodeURL, namespace, name))
	fmt.Println(resp)
	return nil
}

func deleteSpecifiedPod(namespace, name string) error {
	resp := httputil.DeleteWithoutBody(url.Prefix + path.Join(url.PodURL, namespace, name))
	fmt.Println(resp)
	return nil
}

func deleteSpecifiedReplicaSet(namespace, name string) error {
	resp := httputil.DeleteWithoutBody(url.Prefix + path.Join(url.ReplicaSetURL, namespace, name))
	fmt.Println(resp)
	return nil
}

func deleteSpecifiedHPA(namespace, name string) error {
	resp := httputil.DeleteWithoutBody(url.Prefix + path.Join(url.HPAURL, namespace, name))
	fmt.Println(resp)
	return nil
}

func deleteSpecifiedService(namespace, name string) error {
	resp := httputil.DeleteWithoutBody(url.Prefix + path.Join(url.ServiceURL, namespace, name))
	fmt.Println(resp)
	return nil
}

func deleteSpecifiedDNS(namespace, name string) error {
	resp := httputil.DeleteWithoutBody(url.Prefix + path.Join(url.DNSURL, namespace, name))
	fmt.Println(resp)
	return nil
}

func deleteSpecifiedGpuJob(namespace, name string) error {
	resp := httputil.DeleteWithoutBody(url.Prefix + path.Join(url.GpuURL, namespace, name))
	fmt.Println(resp)
	return nil
}

func del(cmd *cobra.Command, args []string) {
	apiObjectType := args[0]
	target := args[1]
	namespace, name := parseName(target)
	apiObjectType = strings.ToLower(apiObjectType)
	var err error
	switch apiObjectType {
	case "node":
		err = deleteSpecifiedNode(namespace, name)
	case "pod":
		err = deleteSpecifiedPod(namespace, name)
	case "rs":
		err = deleteSpecifiedReplicaSet(namespace, name)
	case "hpa":
		err = deleteSpecifiedHPA(namespace, name)
	case "service":
		err = deleteSpecifiedService(namespace, name)
	case "dns":
		err = deleteSpecifiedDNS(namespace, name)
	case "gpu":
		err = deleteSpecifiedGpuJob(namespace, name)
	default:
		err = fmt.Errorf("invalid api object type \"%s\", acceptable api object type is pod, service, etc", apiObjectType)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Printf("Get information of a %s named %s\n", apiObjectType, name)
}
