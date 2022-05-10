package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	cpuPercent  float64
	memPercent  float64
	minReplicas int
	maxReplicas int
)

var autoscaleCmd = &cobra.Command{
	Use:   "autoscale",
	Short: "Kubectl autoscale is used to horizontally and automatically scale a replicaSet according to given metrics.",
	Long:  "Kubectl autoscale is used to horizontally and automatically scale a replicaSet according to given metrics.",
	// two args, target kind(we only support replicaSet now) and its name.
	Args: cobra.ExactArgs(2),
	Run:  autoscale,
}

func configHPA(targetName string) {
	fmt.Printf("Config hpa for target %s\n", targetName)
	fmt.Printf("Metrics: cpu %v, memory %v\n", cpuPercent, memPercent)
	fmt.Printf("Replicas: min %d, max %d\n", minReplicas, maxReplicas)
}

func autoscale(cmd *cobra.Command, args []string) {
	targetKind := args[0]
	if targetKind != "replicaSet" {
		fmt.Println("Unsupported target type!")
		return
	}

	targetName := args[1]
	configHPA(targetName)
}
