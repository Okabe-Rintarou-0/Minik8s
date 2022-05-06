package service

import "testing"

func TestInit(t *testing.T) {
	sm, err := New()
	if err != nil {
		t.Error(err)
	}
	if err = sm.Init(); err != nil {
		t.Error(err)
	}
}

func TestServiceCreate(t *testing.T) {
	sm, err := New()
	if err != nil {
		t.Error(err)
	}
	var eps = make([]EndPoint, 1)
	eps[0] = EndPoint{Name: "KUBE-SEP", Ip: "127.0.0.1", Port: 23333}
	if err = sm.CreateService("KUBE-SVC", "10.96.1.1/32", 32222); err != nil {
		t.Error(err)
	}
	if err = sm.CreateEndpoints("KUBE-SVC", eps); err != nil {
		t.Error(err)
	}
}

func TestServiceDelete(t *testing.T) {
	sm, err := New()
	if err != nil {
		t.Error(err)
	}
	var eps = make([]EndPoint, 1)
	eps[0] = EndPoint{Name: "KUBE-SEP", Ip: "127.0.0.1", Port: 23333}
	if err = sm.DeleteEndPoints("KUBE-SVC", eps); err != nil {
		t.Error(err)
	}
	if err = sm.DeleteService("KUBE-SVC", "10.96.1.1/32", 32222); err != nil {
		t.Error(err)
	}
}

func TestReplicaServiceCreate(t *testing.T) {
	sm, err := New()
	if err != nil {
		t.Error(err)
	}
	var eps = make([]EndPoint, 3)
	eps[0] = EndPoint{Name: "KUBE-SEP1", Ip: "127.0.0.1", Port: 23333}
	eps[1] = EndPoint{Name: "KUBE-SEP2", Ip: "127.0.0.1", Port: 23334}
	eps[2] = EndPoint{Name: "KUBE-SEP3", Ip: "127.0.0.1", Port: 23335}
	if err = sm.CreateService("KUBE-SVC", "10.96.1.1/32", 32222); err != nil {
		t.Error(err)
	}
	if err = sm.CreateEndpoints("KUBE-SVC", eps); err != nil {
		t.Error(err)
	}
}

func TestReplicaServiceDelete(t *testing.T) {
	sm, err := New()
	if err != nil {
		t.Error(err)
	}
	var eps = make([]EndPoint, 3)
	eps[0] = EndPoint{Name: "KUBE-SEP1", Ip: "127.0.0.1", Port: 23333}
	eps[1] = EndPoint{Name: "KUBE-SEP2", Ip: "127.0.0.1", Port: 23334}
	eps[2] = EndPoint{Name: "KUBE-SEP3", Ip: "127.0.0.1", Port: 23335}
	if err = sm.DeleteEndPoints("KUBE-SVC", eps); err != nil {
		t.Error(err)
	}
	if err = sm.DeleteService("KUBE-SVC", "10.96.1.1/32", 32222); err != nil {
		t.Error(err)
	}
}
