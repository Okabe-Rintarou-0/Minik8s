package main

import (
	"fmt"
	"os"
)

func main() {
	template := `# my global config
global:
  scrape_interval:     15s
  evaluation_interval: 15s
alerting:
  alertmanagers:
  - static_configs:
    - targets:

rule_files:

scrape_configs:
  - job_name: 'minik8s'
    static_configs:
    - targets: `
	targets := "['localhost:9090'"
	for _, nodeIp := range os.Args[1:] {
		endpoint := fmt.Sprintf("'%s:8080'", nodeIp)
		targets += ", " + endpoint
	}
	targets += "]\n"

	fmt.Println(template + targets)
}
