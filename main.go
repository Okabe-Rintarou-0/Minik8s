package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testDocker/src/apiObject"
)

func main() {
	content, err := ioutil.ReadFile("./testPod.yaml")

	pod := apiObject.Pod{}
	if err == nil {
		err := yaml.Unmarshal(content, &pod)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("%v\n", pod)
	} else {
		fmt.Println(err.Error())
	}
}
