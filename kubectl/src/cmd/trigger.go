package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/httputil"
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
	URL := url.HttpScheme + url.Hostname + ":8081/" + funcName
	fmt.Println("URL:", URL)
	var errMsg string
	var result string
	resp, err := httputil.PostString(URL, functionData)
	msg := entity.FunctionMsg{}
	if err != nil {
		errMsg = err.Error()
		goto ret
	}
	if err = httputil.ReadAndUnmarshal(resp.Body, &msg); err != nil {
		errMsg = err.Error()
		goto ret
	}

	errMsg = msg.Error
	result = string(msg.Data)
ret:
	fmt.Printf("Called function %s:\n\terror: %s\n\tresult: %s\n", funcName, errMsg, result)
}
