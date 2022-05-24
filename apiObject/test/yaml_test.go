package test

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"minik8s/apiObject"
	"os"
	"testing"
)

func TestYAML(t *testing.T) {
	f, _ := os.Open("../examples/gpu-job/gpu-job-example.yaml")
	content, _ := ioutil.ReadAll(f)
	fmt.Println(string(content))

	gpu := apiObject.GpuJob{}
	yaml.Unmarshal(content, &gpu)
	fmt.Printf("%+v\n", gpu)
}
