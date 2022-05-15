package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/apiserver/src/url"
	"minik8s/util/httputil"
	"strings"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Kubectl delete is used to delete api object given its name",
	Long:  "Kubectl delete is used to delete api object given its name",
	Args:  cobra.ExactValidArgs(2),
	Run:   del,
}

func deleteSpecifiedPod(name string) error {
	resp := httputil.DeleteWithoutBody(url.Prefix + url.PodURL + name)
	fmt.Println(resp)
	return nil
}

func del(cmd *cobra.Command, args []string) {
	apiObjectType := args[0]
	name := args[1]
	apiObjectType = strings.ToLower(apiObjectType)
	var err error
	switch apiObjectType {
	case "pod":
		err = deleteSpecifiedPod(name)
	default:
		err = fmt.Errorf("invalid api object type \"%s\", acceptable api object type is pod, service, etc", apiObjectType)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Printf("Get information of a %s named %s\n", apiObjectType, name)
}
