package helper

import (
	"encoding/json"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"path"
)

func GetNodeHostnames() (hostnames []string) {
	nodeRaws, err := etcd.GetAll(url.NodeURL)
	if err != nil {
		return nil
	}

	for _, nodeRaw := range nodeRaws {
		node := &apiObject.Node{}
		if err = json.Unmarshal([]byte(nodeRaw), node); err == nil {
			hostnames = append(hostnames, node.Name())
		}
	}
	return
}

func ExistsPod(namespace, name string) bool {
	hostnames := GetNodeHostnames()
	var etcdURL string
	for _, hostname := range hostnames {
		etcdURL = path.Join(url.PodURL, hostname, namespace, name)
		if podRaw, err := etcd.Get(etcdURL); err == nil && podRaw != "" {
			return true
		}
	}
	return false
}

func GetPodsApiObjectFromEtcd(node string) (pods []*apiObject.Pod) {
	etcdURL := path.Join(url.PodURL, node)
	if raws, err := etcd.GetAll(etcdURL); err == nil {
		for _, raw := range raws {
			pod := &apiObject.Pod{}
			if err = json.Unmarshal([]byte(raw), &pod); err == nil {
				pods = append(pods, pod)
			}
		}
	}
	return
}
