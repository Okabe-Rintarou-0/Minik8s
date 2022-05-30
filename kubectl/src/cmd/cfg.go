package cmd

import (
	"github.com/spf13/cobra"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"strings"
)

var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "Kubectl cfg is used to change global config",
	Long:  "Kubectl cfg is used to change global config",
	Args:  cobra.ExactArgs(1),
	Run:   config,
}

func config(cmd *cobra.Command, args []string) {
	cfg := args[0]
	parts := strings.Split(cfg, "=")
	if len(parts) < 2 {
		return
	}
	key, value := parts[0], parts[1]

	switch key {
	case "sched":
		listwatch.Publish(topicutil.ScheduleStrategyTopic(), value)
	}
}
