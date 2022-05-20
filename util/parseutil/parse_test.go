package parseutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseHPA(t *testing.T) {
	file, _ := os.Open("../../apiObject/cuda/hpa/hpa-example.yaml")
	content, _ := ioutil.ReadAll(file)
	hpa, err := ParseHPA(content)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Got hpa: %v, metrics are %v\n", hpa, hpa.Metrics())
}
