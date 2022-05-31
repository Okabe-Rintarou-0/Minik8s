package test

import (
	"fmt"
	"minik8s/serverless/src/function"
	"testing"
)

func TestRegistry(t *testing.T){
	err:=function.RemoveFunctionImage("hello1")
	if err!=nil {
		fmt.Println(err)
	}
}
