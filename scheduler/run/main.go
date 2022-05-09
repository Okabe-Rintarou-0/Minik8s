package main

import "minik8s/scheduler/src/scheduler"

func main() {
	s := scheduler.New()
	s.Start()
}
