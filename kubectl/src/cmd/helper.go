package cmd

import (
	"fmt"
	"strings"
)

func parseName(target string) (namespace, name string) {
	parts := strings.Split(target, "/")
	fmt.Println("part", parts)
	if len(parts) == 1 {
		return "default", parts[0]
	} else {
		return parts[0], parts[1]
	}
}
