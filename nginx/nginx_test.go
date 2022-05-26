package nginx

import (
	"testing"
)

const UID = "e466604a-6210-4277-9718-0dc6e4a7ce5f"

func TestApply(t *testing.T) {

	nm := New(UID)
	servers := make([]Server, 1)

	servers[0].Port = 80
	servers[0].Locations = make([]Location, 2)
	servers[0].Locations[0].Addr = "/baidu/"
	servers[0].Locations[0].Dest = "www.baidu.com"
	servers[0].Locations[1].Addr = "/nginx/docs/"
	servers[0].Locations[1].Dest = "nginx.org/en/docs"

	if err := nm.Apply(servers); err != nil {
		t.Error(err)
	}
}

func TestStart(t *testing.T) {
	if err := New(UID).Start(); err != nil {
		t.Error(err)
	}
}

func TestShutdown(t *testing.T) {
	if err := New(UID).Shutdown(); err != nil {
		t.Error(err)
	}
}
