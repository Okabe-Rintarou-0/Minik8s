package test

import (
	"encoding/json"
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

func TestJson(t *testing.T) {
	f, _ := os.Open("../examples/workflow/workflow.json")
	content, _ := ioutil.ReadAll(f)
	fmt.Println(string(content))

	wf := apiObject.Workflow{}
	json.Unmarshal(content, &wf)
	fmt.Printf("%+v\n", wf)

	mmap := map[string]interface{}{
		"1": 1,
	}
	res, _ := json.Marshal(mmap)
	fmt.Println(string(res))

	json.Unmarshal(res, &mmap)
	fmt.Println(mmap["1"].(float64) == float64(1))
}
