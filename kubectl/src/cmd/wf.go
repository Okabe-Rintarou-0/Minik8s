package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/util/httputil"
)

var wfCmd = &cobra.Command{
	Use:   "wf",
	Short: "Kubectl wf is used to apply a workflow",
	Long:  "Kubectl wf is used to apply a workflow",
	Args:  cobra.ExactArgs(1),
	Run:   handleWorkflow,
}

func addWorkflowToApiServer(wf *apiObject.Workflow) error {
	URL := url.Prefix + url.WorkflowURL
	if resp, err := httputil.PostJson(URL, wf); err == nil {
		var content []byte
		if content, err = ioutil.ReadAll(resp.Body); err == nil {
			defer resp.Body.Close()
			fmt.Println(string(content))
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func handleWorkflow(cmd *cobra.Command, args []string) {
	op := args[0]
	var err error
	switch op {
	case "apply":
		var raw []byte
		raw, err = ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		wf := apiObject.Workflow{}
		if err = json.Unmarshal(raw, &wf); err != nil {
			fmt.Println(err.Error())
			return
		}

		err = addWorkflowToApiServer(&wf)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}
