package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}

var filePath string

func init() {
	applyCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "filePath of api object yaml file")

	autoscaleCmd.Flags().StringVarP(&target, "target", "t", "", "target name")
	autoscaleCmd.Flags().Float64VarP(&cpuPercent, "cpu", "c", 0.0, "cpu utilization percent metric")
	autoscaleCmd.Flags().Float64VarP(&memPercent, "mem", "m", 0.0, "memory utilization percent metric")
	autoscaleCmd.Flags().IntVarP(&minReplicas, "min", "", 1, "min replicas")
	autoscaleCmd.Flags().IntVarP(&maxReplicas, "max", "", 1, "max replicas")

	labelCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "overwrite labels")

	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(describeCmd)
	rootCmd.AddCommand(autoscaleCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(labelCmd)
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
