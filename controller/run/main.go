package main

import (
	"minik8s/controller/src/controller"
)

func main() {
	cm := controller.NewControllerManager()
	cm.Start()
}
