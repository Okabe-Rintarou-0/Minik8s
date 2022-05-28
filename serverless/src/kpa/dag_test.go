package kpa

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"minik8s/apiObject"
	"os"
	"testing"
)

func TestDAG(t *testing.T) {
	f, _ := os.Open("./hello.json")
	content, _ := ioutil.ReadAll(f)
	wf := apiObject.Workflow{}
	_ = json.Unmarshal(content, &wf)
	fmt.Printf("%+v\n", wf)
	dag := Workflow2DAG(&wf)
	if dag != nil {
		TraverseDAG(dag)
		fmt.Println(gjson.Get(dag.EntryParams, "a").Int())
		fmt.Println(gjson.Get(dag.EntryParams, "b").Int())
		fmt.Println(gjson.Get(dag.EntryParams, "name").String())
	}
}
