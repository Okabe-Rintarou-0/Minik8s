package test

import (
	"log"
	"minik8s/apiserver/src/etcd"
	"testing"
)

func TestEtcd(t *testing.T) {
	etcd.Put("22", "33")
	etcd.Put("22", "1234")
	log.Print("22:", etcd.Get("22"))
	log.Print("33:", etcd.Get("33"))

	etcd.Delete("22")
	log.Print("22:", etcd.Get("22"))
}
