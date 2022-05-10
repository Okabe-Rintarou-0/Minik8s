package main

import "minik8s/kubelet/src/kubelet"

func main() {
	kl := kubelet.New()
	kl.Run()
}
