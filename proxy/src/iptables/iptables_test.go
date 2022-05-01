package iptables

import (
	"testing"
)

func TestServiceCreate(t *testing.T) {
	im, err := New()
	if err != nil {
		t.Error(err)
	}
	if err = im.NewService("KUBE-SERVICES"); err != nil {
		t.Error(err)
	}
	if err = im.NewChain("KUBE-SEP"); err != nil {
		t.Error(err)
	}
	if err = im.NewChain("KUBE-SVC"); err != nil {
		t.Error(err)
	}
	if err = im.AppendChainToEndpoint("KUBE-SEP", "127.0.0.1:23333"); err != nil {
		t.Error(err)
	}
	if err = im.AppendChainToChains("KUBE-SVC", []string{"KUBE-SEP"}); err != nil {
		t.Error(err)
	}
	if err = im.AppendServiceToChain("KUBE-SERVICES", "KUBE-SVC", "10.96.1.1/32", 32222); err != nil {
		t.Error(err)
	}
}

func TestReplicaServiceCreate(t *testing.T) {
	im, err := New()
	if err != nil {
		t.Error(err)
	}
	if err = im.NewService("KUBE-SERVICES"); err != nil {
		t.Error(err)
	}
	if err = im.NewChain("KUBE-SEP3"); err != nil {
		t.Error(err)
	}
	if err = im.NewChain("KUBE-SEP4"); err != nil {
		t.Error(err)
	}
	if err = im.NewChain("KUBE-SEP5"); err != nil {
		t.Error(err)
	}
	if err = im.NewChain("KUBE-SVC"); err != nil {
		t.Error(err)
	}
	if err = im.AppendChainToEndpoint("KUBE-SEP3", "127.0.0.1:23333"); err != nil {
		t.Error(err)
	}
	if err = im.AppendChainToEndpoint("KUBE-SEP4", "127.0.0.1:23334"); err != nil {
		t.Error(err)
	}
	if err = im.AppendChainToEndpoint("KUBE-SEP5", "127.0.0.1:23335"); err != nil {
		t.Error(err)
	}
	if err = im.AppendChainToChains("KUBE-SVC", []string{"KUBE-SEP3", "KUBE-SEP4", "KUBE-SEP5"}); err != nil {
		t.Error(err)
	}
	if err = im.AppendServiceToChain("KUBE-SERVICES", "KUBE-SVC", "10.96.1.1/32", 32222); err != nil {
		t.Error(err)
	}
}
