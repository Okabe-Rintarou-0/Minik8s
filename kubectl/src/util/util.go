package util

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"minik8s/apiObject"
	"os"
	"strings"
)

type ApiObjectType byte

const (
	Unknown ApiObjectType = iota
	Node
	Pod
	Service
	Deployment
)

func (tp *ApiObjectType) String() string {
	switch *tp {
	case Pod:
		return "Pod"
	case Node:
		return "Node"
	case Service:
		return "Service"
	case Deployment:
		return "Deployment"
	}
	return "Unknown"
}

func IsValidApiObjectType(objectType string) bool {
	return objectType == "pod" || objectType == "pods" ||
		objectType == "deployment" || objectType == "service"
}

func isLetter(char rune) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'
}

func parseType(content []byte) ApiObjectType {
	idx := strings.LastIndex(string(content), "kind:") + 5
	total := len(content)
	// Ignore spaces
	for idx < total && content[idx] == ' ' {
		idx++
	}

	startIdx := idx
	endIdx := idx
	for endIdx < total && isLetter(rune(content[endIdx])) {
		endIdx++
	}
	kind := string(content[startIdx:endIdx])
	switch kind {
	case "Pod":
		return Pod
	}
	return Unknown
}

func LoadContent(filePath string) ([]byte, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Read content: %s \nfrom %s\n", content, filePath)
	return content, nil
}

func ParseApiObjectType(content []byte) (ApiObjectType, error) {
	tp := parseType(content)

	fmt.Println("Api object's type is", tp.String())

	return tp, nil
}

func ParsePod(content []byte) (*apiObject.Pod, error) {
	pod := &apiObject.Pod{}
	err := yaml.Unmarshal(content, pod)
	return pod, err
}
