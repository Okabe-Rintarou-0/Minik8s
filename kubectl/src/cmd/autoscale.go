package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/apiserver/src/url"
	"minik8s/util/httputil"
	"strconv"
)

var (
	target      string
	cpuPercent  float64
	memPercent  float64
	minReplicas int
	maxReplicas int
)

var autoscaleCmd = &cobra.Command{
	Use:   "autoscale",
	Short: "Kubectl autoscale is used to horizontally and automatically scale a replicaSet according to given metrics.",
	Long:  "Kubectl autoscale is used to horizontally and automatically scale a replicaSet according to given metrics.",
	// two args, target name.
	Args: cobra.ExactArgs(1),
	Run:  autoscale,
}

func configHPAToApiServer(name string) {
	fmt.Printf("Config hpa name = %s\n", name)
	fmt.Printf("Target: %v\n", target)
	fmt.Printf("Metrics: cpu %v, memory %v\n", cpuPercent, memPercent)
	fmt.Printf("Replicas: min %d, max %d\n", minReplicas, maxReplicas)

	httputil.PostForm(url.Prefix+url.AutoscaleURL+name, map[string]string{
		"target": target,
		"cpu":    strconv.FormatFloat(cpuPercent, 'f', 4, 64),
		"mem":    strconv.FormatFloat(memPercent, 'f', 4, 64),
		"min":    strconv.FormatInt(int64(minReplicas), 10),
		"max":    strconv.FormatInt(int64(maxReplicas), 10),
	})
}

func autoscale(cmd *cobra.Command, args []string) {
	name := args[0]
	configHPAToApiServer(name)
}
