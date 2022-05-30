package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/util/httputil"
	"path"
	"strings"
)

var funcCmd = &cobra.Command{
	Use:   "func",
	Short: "Kubectl func is used to apply a func",
	Long:  "Kubectl func is used to apply a func",
	Args:  cobra.ExactArgs(1),
	Run:   handleFunc,
}

var (
	function     string
	functionPath string
)

func updateFuncToApiServer(function *apiObject.Function) {
	URL := url.Prefix + path.Join(url.FuncURL, function.Name)
	if resp, err := httputil.PutJson(URL, function); err != nil {
		fmt.Println(err.Error())
	} else if content, err := ioutil.ReadAll(resp.Body); err == nil {
		defer resp.Body.Close()
		fmt.Println(string(content))
	}
}

func addFuncToApiServer(function *apiObject.Function) {
	URL := url.Prefix + url.FuncURL
	if resp, err := httputil.PostJson(URL, function); err != nil {
		fmt.Println(err.Error())
	} else if content, err := ioutil.ReadAll(resp.Body); err == nil {
		defer resp.Body.Close()
		fmt.Println(string(content))
	}
}

func removeFuncToApiServer(function string) {
	URL := url.Prefix + path.Join(url.FuncURL, function)
	resp := httputil.DeleteWithoutBody(URL)
	fmt.Println(resp)
}

func handleFunc(cmd *cobra.Command, args []string) {
	op := args[0]
	if function == "" {
		fmt.Println("Function can not be empty!")
		return
	}

	op = strings.ToLower(op)
	switch op {
	case "add":
		if functionPath == "" {
			fmt.Println("Function path can not be empty!")
			return
		}
		addFuncToApiServer(&apiObject.Function{
			Name: function,
			Path: functionPath,
		})
	case "rm":
		removeFuncToApiServer(function)
	case "update":
		if functionPath == "" {
			fmt.Println("Function path can not be empty!")
			return
		}
		updateFuncToApiServer(&apiObject.Function{
			Name: function,
			Path: functionPath,
		})
	}
}
