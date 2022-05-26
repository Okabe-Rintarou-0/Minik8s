package main

import (
	"fmt"
	"minik8s/serverless/src/function"
	"minik8s/serverless/src/knative"
	"minik8s/serverless/src/registry"
	"os/exec"
)

func main() {
	registry.InitRegistry()
	// only run on master node
	cmd := exec.Command("ls")
	output, _ := cmd.Output()
	fmt.Println(string(output))
	function.InitFunction("helloworld", "../src/app/func.py") // the third parameter need to be replaced
	kn := knative.NewKnative()
	kn.Run()
}
