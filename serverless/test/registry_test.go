package test

import (
	"minik8s/serverless/src/function"
	"minik8s/serverless/src/registry"
	"testing"
)

func TestRegistry(t *testing.T) {
	registry.InitRegistry()
	function.InitFunction()
}

