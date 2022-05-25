package utils

import (
	"minik8s/kubelet/src/runtime/container"
	"minik8s/kubelet/src/runtime/image"
	"log"
)


var (
	cm = container.NewContainerManager()
	im = image.NewImageManager()
)

// pull image from web
func PullImg(imageName string) {
	if exists, err := im.ExistsImage(imageName); !exists && err == nil {
		err := im.PullImage(imageName, &image.PullConfig{All: false, Verbose: true})
		log.Printf("pull image of %s complete\n", imageName)
		if err != nil {
			log.Print(err.Error())
			return
		}
	} else if err == nil {
		log.Printf("image of %s exists", imageName)
	} else {
		log.Print(err.Error())
		return
	}
}
