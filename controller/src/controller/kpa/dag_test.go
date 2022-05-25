package kpa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minik8s/apiObject"
	"os"
	"testing"
)

func TestDAG(t *testing.T) {
	f, _ := os.Open("./workflow.json")
	content, _ := ioutil.ReadAll(f)
	wf := apiObject.Workflow{}
	_ = json.Unmarshal(content, &wf)
	fmt.Printf("%+v\n", wf)
	dag := Workflow2DAG(&wf)
	if dag != nil {
		TraverseDAG(dag)
	}
}
