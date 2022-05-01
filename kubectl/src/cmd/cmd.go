package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/kubectl/src/util"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}

var filePath string

func init() {
	applyCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "filePath of api object yaml file")
	rootCmd.AddCommand(applyCmd)
}

var rootCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "kubectl is for better control of minik8s",
	Long: `by using kubectl, you can create api object in minik8s, or know details of them by using kubectl describe command.
For example: kubectl apply -f ./example.yaml; kubectl describe pod examplePod`,
	Run: runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {
	// Reach here if no args
	fmt.Printf("execute %s args:%v \n", cmd.Name(), args)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "kubectl apply command",
	Long:  "kubectl apply is used to create api object in a declarative way.",
	Run:   apply,
}

func apply(cmd *cobra.Command, args []string) {
	fmt.Println(*cmd)
	fmt.Println("apply ", args)
	fmt.Println(filePath)
	content, err := util.LoadContent(filePath)
	if err != nil {
		fmt.Println(err.Error())
	}

	tp, err := util.ParseApiObjectType(content)
	if err != nil {
		fmt.Println(err.Error())
	}

	switch tp {
	case util.Pod:
		pod, err := util.ParsePod(content)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Get pod %v\n", *pod)
	}
}
