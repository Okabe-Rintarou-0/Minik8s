package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/kubelet/src/kubelet"
	"minik8s/util/httputil"
	"minik8s/util/netutil"
	"minik8s/util/wait"
	"os"
	"time"
)

func registerNode(ip string) {
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
		Ip: ip,
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
	var ip string
	flag.StringVar(&ip, "ip", "127.0.0.1", "ip address for node register")
	flag.Parse()

	go wait.Period(0, time.Minute, func() {
		registerNode(ip)
	})

	kl := kubelet.New()
	kl.Run()
}
