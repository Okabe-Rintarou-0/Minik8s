package test

import (
	// "io/ioutil"
	// "os"
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

func TestJsonParse(t *testing.T) {
	// f, _ := os.Open("workflow.json")
	// content, _ := ioutil.ReadAll(f)
	// str:=string(content)

	const json = `{
		"name":{
			"first":"Janet",
			"last":"Prichard"
		},
		"world":[2,3],
		"app":23.33,
		"arg":33,
		"age":47
	}`

	// wf := apiObject.Workflow{}
	// json.Unmarshal(content, &wf)
	// fmt.Printf("%+v\n", wf)

	res:=gjson.Get(json,"arg")
	fmt.Print(res.Type,"\n")
}