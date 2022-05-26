package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var funcCmd = &cobra.Command{
	Use:   "func",
	Short: "Kubectl func is used to apply a func",
	Long:  "Kubectl func is used to apply a func",
	Run:   handleFunc,
}

var (
	function     string
	functionPath string
)

func handleFunc(cmd *cobra.Command, args []string) {
	if function == "" {
		fmt.Println("Function can not be empty!")
		return
	}

	if functionPath == "" {
		fmt.Println("Function path can not be empty!")
		return
	}
}
