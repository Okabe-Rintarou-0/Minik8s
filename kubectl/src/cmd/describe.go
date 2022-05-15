package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/httputil"
)

func getPodDescriptionFromApiServer(name string) (desc *entity.PodDescription, err error) {
	err = httputil.GetAndUnmarshal(url.Prefix+url.PodDescriptionURL+name, &desc)
	return
}

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Kubectl describe is used to get detailed information of the api object with given unique name",
	Long:  "Kubectl describe is used to get detailed information of the api object with given unique name",
	// exactly two args, one is the type of api object, another is the unique name
	// for example, kubectl describe pod example, where pod is the type and example is the unique name
	Args: cobra.ExactValidArgs(2),
	Run:  describe,
}

func describe(cmd *cobra.Command, args []string) {
	apiObjectType := args[0]
	name := args[1]

	var err error
	switch apiObjectType {
	case "pod":
		err = printSpecifiedPodDescription(name)
	default:
		fmt.Println("Invalid api object type!")
		return
	}

	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Printf("Describe information of a %s named %s\n", apiObjectType, name)
}
