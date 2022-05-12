package main

import "minik8s/apiserver/src/apiserver"

func main() {
	apiServer := apiserver.New()
	apiServer.Run()
}
