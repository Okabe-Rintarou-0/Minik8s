package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/apiserver/src/url"
	"minik8s/util/httputil"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Kubectl reset is used to reset minik8s persistent status",
	Long:  "Kubectl reset is used to reset minik8s persistent status",
	Run:   reset,
}

func reset(cmd *cobra.Command, args []string) {
	URL := url.Prefix + url.ResetURL
	resp := httputil.DeleteWithoutBody(URL)
	fmt.Println(resp)
}
