package main

import (
	"fmt"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/kubelet/src/kubelet"
	"minik8s/util/httputil"
	"minik8s/util/netutil"
	"os"
)

func registerNode() {
	hostname := netutil.Hostname()
	node := apiObject.Node{
		Base: apiObject.Base{
			ApiVersion: "v1",
			Kind:       "Node",
			Metadata: apiObject.Metadata{
				Name:      hostname,
				Namespace: "default",
				UID:       "",
			},
		},
		Ip: "0.0.0.0",
	}

	URL := url.Prefix + url.NodeURL
	resp, err := httputil.PostJson(URL, node)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	body := resp.Body
	content, _ := ioutil.ReadAll(body)
	defer body.Close()
	fmt.Println(string(content))
}

func main() {
	registerNode()

	kl := kubelet.New()
	kl.Run()
}
