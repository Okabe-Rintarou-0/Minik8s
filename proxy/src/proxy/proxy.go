package proxy

import "minik8s/proxy/src/service"

type Proxy struct {
	iptablesManager *service.Manager
}
