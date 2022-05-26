package utils

import (
	"fmt"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/kubelet/src/runtime/image"
)

var (
	cm = container.NewContainerManager()
	im = image.NewImageManager()
)

// PullImg pull image from web
func PullImg(imageName string) error {
	if exists, err := im.ExistsImage(imageName); !exists && err == nil {
		err := im.PullImage(imageName, &image.PullConfig{All: false, Verbose: true})
		fmt.Printf("pull image of %s complete\n", imageName)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	} else if err == nil {
		fmt.Printf("image of %s exists\n", imageName)
	} else {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
