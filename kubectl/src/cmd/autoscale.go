package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/apiserver/src/url"
	"minik8s/util/httputil"
	"path"
	"strconv"
)

var (
	target        string
	cpuPercent    float64
	memPercent    float64
	minReplicas   int
	maxReplicas   int
	scaleInterval int
)

var autoscaleCmd = &cobra.Command{
	Use:   "autoscale",
	Short: "Kubectl autoscale is used to horizontally and automatically scale a replicaSet according to given metrics.",
	Long:  "Kubectl autoscale is used to horizontally and automatically scale a replicaSet according to given metrics.",
	// two args, target name.
	Args: cobra.ExactArgs(1),
	Run:  autoscale,
}

func configHPAToApiServer(namespace, name string) {
	fmt.Printf("Config hpa %s/%s\n", namespace, name)
	fmt.Printf("Target: %v\n", target)
	fmt.Printf("Metrics: cpu %v, memory %v\n", cpuPercent, memPercent)
	fmt.Printf("Replicas: min %d, max %d\n", minReplicas, maxReplicas)
	fmt.Printf("Interval: %d\n", scaleInterval)

	resp := httputil.PostForm(url.Prefix+path.Join(url.AutoscaleURL, namespace, name), map[string]string{
		"target":   target,
		"cpu":      strconv.FormatFloat(cpuPercent, 'f', 4, 64),
		"mem":      strconv.FormatFloat(memPercent, 'f', 4, 64),
		"min":      strconv.FormatInt(int64(minReplicas), 10),
		"max":      strconv.FormatInt(int64(maxReplicas), 10),
		"interval": strconv.FormatInt(int64(scaleInterval), 10),
	})
	fmt.Println(resp)
}

func autoscale(cmd *cobra.Command, args []string) {
	fullName := args[0]
	namespace, name := parseName(fullName)
	configHPAToApiServer(namespace, name)
}
