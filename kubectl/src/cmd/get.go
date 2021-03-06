package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Kubectl get is used to get brief information of the api object with given unique name",
	Long:  "Kubectl get is used to get brief information of the api object with given unique name",
	// exactly two args, one is the type of api object, another is the unique name
	// for example, kubectl get pod example, where pod is the type and example is the unique name
	Args: cobra.MinimumNArgs(1),
	Run:  get,
}

func get(cmd *cobra.Command, args []string) {
	apiObjectType := args[0]
	var name string
	if len(args) > 1 {
		name = args[1]
	}

	apiObjectType = strings.ToLower(apiObjectType)
	var err error
	switch apiObjectType {
	case "pod":
		err = printSpecifiedPodStatus(name)
	case "pods":
		err = printPodStatuses()
	case "node":
		err = printSpecifiedNodeStatus(name)
	case "nodes":
		err = printNodeStatuses()
	case "rs":
		err = printSpecifiedReplicaSetStatus(name)
	case "rss":
		err = printReplicaSetStatuses()
	case "hpas":
		err = printHPAStatuses()
	case "hpa":
		err = printSpecifiedHPAStatus(name)
	case "service":
		err = printSpecifiedService(name)
	case "services":
		err = printServices()
	case "dns":
		err = printSpecifiedDns(name)
	case "dnses":
		err = printDnses()
	case "wf":
		err = printSpecifiedWorkflow(name)
	case "wfs":
		err = printWorkflows()
	case "func":
		err = printSpecifiedFunction(name)
	case "funcs":
		err = printFunctions()
	case "gpu":
		err = printSpecifiedGpuJob(name)
	case "gpus":
		err = printGpuJobs()
	default:
		err = fmt.Errorf("invalid api object type \"%s\", acceptable api object type is pod, service, etc", apiObjectType)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Printf("Get information of a %s named %s\n", apiObjectType, name)
}
