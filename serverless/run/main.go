package main

import (
	"minik8s/serverless/src/function"
	"minik8s/serverless/src/registry"
)

func main() {
	registry.InitRegistry()
	function.InitFunction()
}
