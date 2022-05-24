package main

import (
	"minik8s/proxy/src/proxy"
	"minik8s/util/logger"
)

var log = logger.Log("Proxy-main")

func main() {
	if p, err := proxy.New(); err != nil {
		log(err.Error())
	} else {
		p.Run()
	}
}
