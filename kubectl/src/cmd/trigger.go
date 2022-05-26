package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/entity"
	functrigger "minik8s/serverless/src/trigger"
)

var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Kubectl trigger is used to trigger serverless function",
	Long:  "Kubectl trigger is used to trigger serverless function",
	Args:  cobra.ExactArgs(1),
	Run:   trigger,
}

var functionData string

func trigger(cmd *cobra.Command, args []string) {
	funcName := args[0]
	result, err := functrigger.Trigger(funcName, entity.FunctionData(functionData))
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	fmt.Printf("Called function %s:\n\terror: %s\n\tresult: %s\n", funcName, errMsg, result)
}
