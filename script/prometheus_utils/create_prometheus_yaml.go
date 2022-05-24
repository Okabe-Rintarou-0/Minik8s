package main

import (
	"fmt"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/httputil"
)

func getNodes() (nodes []string) {
	var nodeStatuses []*entity.NodeStatus
	if err := httputil.GetAndUnmarshal(url.Prefix+url.NodeURL, &nodeStatuses); err == nil {
		for _, status := range nodeStatuses {
			nodes = append(nodes, status.Ip)
		}
	}
	return
}

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
	for _, nodeIp := range getNodes() {
		endpoint := fmt.Sprintf("'%s:8000'", nodeIp)
		targets += ", " + endpoint
	}
	targets += "]\n"

	fmt.Println(template + targets)
}
