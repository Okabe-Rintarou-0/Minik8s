package util

import (
	"io/ioutil"
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
	ReplicaSet
	HorizontalPodAutoscaler
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
	case ReplicaSet:
		return "ReplicaSet"
	case HorizontalPodAutoscaler:
		return "HorizontalPodAutoscaler"
	}
	return "Unknown"
}

func isLetter(char rune) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'
}

func parseType(content []byte) ApiObjectType {
	idx := strings.Index(string(content), "kind:") + 5
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
	case "ReplicaSet":
		return ReplicaSet
	case "HorizontalPodAutoscaler":
		return HorizontalPodAutoscaler
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
	//fmt.Printf("Read content: %s \nfrom %s\n", content, filePath)
	return content, nil
}

func ParseApiObjectType(content []byte) (ApiObjectType, error) {
	tp := parseType(content)
	//fmt.Println("Api object's type is", tp.String())
	return tp, nil
}
