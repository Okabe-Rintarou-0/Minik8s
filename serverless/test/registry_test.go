package test

import (
	"minik8s/serverless/src/registry"
	"testing"
)

func TestRegistry(t *testing.T) {
	registry.InitRegistry("127.0.0.1")
	registry.PushImage("python:3.10-slim")

}

//func TestFunction(t *testing.T) {
//	function.InitFunction()
//	utils.CreateContainer("test1", "python:3.10-slim")
//}
