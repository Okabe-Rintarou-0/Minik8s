package utils

import (
	"log"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/kubelet/src/runtime/image"
)

//const etcdImage = "bitnami/etcd:3.5"
//const serverName = "etcd-server"

var (
	cm = container.NewContainerManager()
	im = image.NewImageManager()
)

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

//func Start() {
//	pullImage()
//	ID, State := findContainer()
//	if ID != "" {
//		log.Printf("[etcd] Find etcd container, ID: %s, State: %s\n", ID, State)
//
//	} else {
//		ID = createContainer()
//		log.Printf("[etcd] Create Container, ID: %s\n", ID)
//	}
//	if State != "running" && ID != "" {
//		startContainer(ID)
//		log.Printf("[etcd] Container starts\n")
//	}
//}
