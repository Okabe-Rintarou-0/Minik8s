package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/kubectl/src/util"
	"minik8s/util/parseutil"
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
	case util.Pod:
		pod, err := parseutil.ParsePod(content)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Get pod %v\n", *pod)
	}
}
