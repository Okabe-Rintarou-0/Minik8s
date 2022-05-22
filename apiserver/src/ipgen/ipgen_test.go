package ipgen

import (
	"minik8s/apiserver/src/url"
	"testing"
)

func Test(t *testing.T) {
	ig, err := New(url.PodIpGeneratorURL, "10.44.0.0/28")
	if err != nil {
		t.Error(err)
	}
	if err := ig.Clear(); err != nil {
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
