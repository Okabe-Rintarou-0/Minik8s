package run

import "minik8s/serverless/src/registry"

func main() {
	registry.InitRegistry("127.0.0.1")
}
