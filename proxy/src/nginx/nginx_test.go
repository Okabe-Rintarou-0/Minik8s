package nginx

import (
	"testing"
)

func Test(t *testing.T) {
	nm := New("/etc/nginx/nginx.conf")
	servers := make([]Server, 2)

	servers[0].Port = 8010
	servers[0].Locations = make([]Location, 2)
	servers[0].Locations[0].Addr = "/baidu/"
	servers[0].Locations[0].Dest = "http://www.baidu.com/"
	servers[0].Locations[1].Addr = "/nginx/docs/"
	servers[0].Locations[1].Dest = "http://nginx.org/en/docs/"

	servers[1].Port = 8011
	servers[1].Locations = make([]Location, 2)
	servers[1].Locations[0].Addr = "/kubernetes/"
	servers[1].Locations[0].Dest = "http://kubernetes.io/"
	servers[1].Locations[1].Addr = "/coredns/"
	servers[1].Locations[1].Dest = "http://coredns.io/"

	if err := nm.Apply(servers); err != nil {
		t.Error(err)
	}
}
