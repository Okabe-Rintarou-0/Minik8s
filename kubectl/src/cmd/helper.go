package cmd

import (
	"strings"
)

func parseName(target string) (namespace, name string) {
	parts := strings.Split(target, "/")
	if len(parts) == 1 {
		return "default", parts[0]
	} else {
		return parts[0], parts[1]
	}
}
