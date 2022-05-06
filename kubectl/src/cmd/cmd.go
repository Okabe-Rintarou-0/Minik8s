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
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(describeCmd)
}

var rootCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "Kubectl is for better control of minik8s",
	Long: `By using kubectl, you can create api object in minik8s, or know details of them by using kubectl describe command.
For example: kubectl apply -f ./example.yaml; kubectl describe pod examplePod`,
	Run: runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {
	// Reach here if no args
	fmt.Printf("execute %s args:%v \n", cmd.Name(), args)
}

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
		pod, err := util.ParsePod(content)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Get pod %v\n", *pod)
	}
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Kubectl get is used to get brief information of the api object with given unique name",
	Long:  "Kubectl get is used to get brief information of the api object with given unique name",
	// exactly two args, one is the type of api object, another is the unique name
	// for example, kubectl get pod example, where pod is the type and example is the unique name
	Args: cobra.ExactValidArgs(2),
	Run:  get,
}

func get(cmd *cobra.Command, args []string) {
	apiObjectType := args[0]
	name := args[1]

	if !util.IsValidApiObjectType(apiObjectType) {
		fmt.Printf("invalid api object type \"%s\", acceptable api object type is pod, service, etc.", apiObjectType)
		return
	}

	fmt.Printf("Get information of a %s named %s\n", apiObjectType, name)
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

	if !util.IsValidApiObjectType(apiObjectType) {
		fmt.Println("Invalid api object type!")
		return
	}

	fmt.Printf("Describe information of a %s named %s\n", apiObjectType, name)
}
