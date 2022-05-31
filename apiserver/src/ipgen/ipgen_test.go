package ipgen

import (
	"minik8s/apiserver/src/url"
	"testing"
)

func Test(t *testing.T) {
	ig := New(url.PodIpURL, 16)
	if err := ig.Clear(url.PodIpBase); err != nil {
		t.Error(err)
	}
	t.Log(ig.GetCurrent())
	t.Log(ig.GetNext())
	t.Log(ig.GetCurrent())
	t.Log(ig.GetCurrent())
	t.Log(ig.GetCurrentWithMask())
	t.Log(ig.GetNext())
	t.Log(ig.GetNext())
	for i := 1; i < 200; i++ {
		t.Log(ig.GetNext())
		t.Log(ig.GetNextWithMask())
	}
}
