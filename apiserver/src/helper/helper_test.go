package helper

import (
	"encoding/json"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"minik8s/util/uidutil"
	"path"
	"testing"
)

func getPod1() apiObject.Pod {
	pod := apiObject.Pod{}
	pod.Metadata.UID = uidutil.New()
	pod.AddLabel("key1", "value1")
	pod.AddLabel("key2", "value2")
	pod.Spec.Containers = make([]apiObject.Container, 2)
	pod.Spec.Containers[0].Ports = make([]apiObject.ContainerPort, 2)
	pod.Spec.Containers[0].Ports[0].HostIP = "10.44.0.1"
	pod.Spec.Containers[0].Ports[0].HostPort = "80"
	pod.Spec.Containers[0].Ports[1].HostIP = "10.44.0.2"
	pod.Spec.Containers[0].Ports[1].HostPort = "80"
	pod.Spec.Containers[1].Ports = make([]apiObject.ContainerPort, 2)
	pod.Spec.Containers[1].Ports[0].HostIP = "10.44.1.1"
	pod.Spec.Containers[1].Ports[0].HostPort = "81"
	pod.Spec.Containers[1].Ports[1].HostIP = "10.44.1.2"
	pod.Spec.Containers[1].Ports[1].HostPort = "81"
	return pod
}

func getPod2() apiObject.Pod {
	pod := apiObject.Pod{}
	pod.Metadata.UID = uidutil.New()
	pod.AddLabel("key1", "value1")
	pod.AddLabel("key3", "value3")
	pod.Spec.Containers = make([]apiObject.Container, 2)
	pod.Spec.Containers[0].Ports = make([]apiObject.ContainerPort, 2)
	pod.Spec.Containers[0].Ports[0].HostIP = "11.44.0.1"
	pod.Spec.Containers[0].Ports[0].HostPort = "80"
	pod.Spec.Containers[0].Ports[1].HostIP = "11.44.0.2"
	pod.Spec.Containers[0].Ports[1].HostPort = "80"
	pod.Spec.Containers[1].Ports = make([]apiObject.ContainerPort, 2)
	pod.Spec.Containers[1].Ports[0].HostIP = "11.44.1.1"
	pod.Spec.Containers[1].Ports[0].HostPort = "81"
	pod.Spec.Containers[1].Ports[1].HostIP = "11.44.1.2"
	pod.Spec.Containers[1].Ports[1].HostPort = "81"
	return pod
}

func Test(t *testing.T) {
	var err error

	if err = etcd.DeleteAllKeys(); err != nil {
		t.Error(err)
	}

	pod1 := getPod1()
	if err = AddEndpoints(pod1); err != nil {
		t.Error(err)
	}
	t.Log(pod1.Metadata.UID)
	pod2 := getPod2()
	if err = AddEndpoints(pod2); err != nil {
		t.Error(err)
	}
	t.Log(pod2.Metadata.UID)

	//if err = DelEndpoints(pod1); err != nil {
	//	t.Error(err)
	//}
	//if err = DelEndpoints(pod2); err != nil {
	//	t.Error(err)
	//}

	var str string
	var endpoints etcd.Endpoints

	if str, err = etcd.Get(path.Join(url.EndpointURL, "key1", "value1")); err != nil {
		t.Error(err)
	} else {
		t.Log(str)
	}
	if err := json.Unmarshal([]byte(str), &endpoints); err != nil {
		t.Error(err)
	}
	for _, UID := range endpoints.Get(pod1.Metadata.UID) {
		if str, err = etcd.Get(path.Join(url.EndpointURL, UID)); err != nil {
			t.Error(err)
		} else {
			t.Log(str)
		}
	}

	if str, err = etcd.Get(path.Join(url.EndpointURL, "key2", "value2")); err != nil {
		t.Error(err)
	} else {
		t.Log(str)
	}
	if err := json.Unmarshal([]byte(str), &endpoints); err != nil {
		t.Error(err)
	}
	for _, UID := range endpoints.Get(pod1.Metadata.UID) {
		if str, err = etcd.Get(path.Join(url.EndpointURL, UID)); err != nil {
			t.Error(err)
		} else {
			t.Log(str)
		}
	}

	if str, err = etcd.Get(path.Join(url.EndpointURL, "key1", "value1")); err != nil {
		t.Error(err)
	} else {
		t.Log(str)
	}
	if err := json.Unmarshal([]byte(str), &endpoints); err != nil {
		t.Error(err)
	}
	for _, UID := range endpoints.Get(pod2.Metadata.UID) {
		if str, err = etcd.Get(path.Join(url.EndpointURL, UID)); err != nil {
			t.Error(err)
		} else {
			t.Log(str)
		}
	}

	if str, err = etcd.Get(path.Join(url.EndpointURL, "key3", "value3")); err != nil {
		t.Error(err)
	} else {
		t.Log(str)
	}
	if err := json.Unmarshal([]byte(str), &endpoints); err != nil {
		t.Error(err)
	}
	for _, UID := range endpoints.Get(pod2.Metadata.UID) {
		if str, err = etcd.Get(path.Join(url.EndpointURL, UID)); err != nil {
			t.Error(err)
		} else {
			t.Log(str)
		}
	}

	if ret, err := getEndpoints("key1", "value1"); err != nil {
		t.Error(err)
	} else {
		t.Log(ret)
	}
	if ret, err := getEndpoints("key2", "value2"); err != nil {
		t.Error(err)
	} else {
		t.Log(ret)
	}
	if ret, err := getEndpoints("key3", "value3"); err != nil {
		t.Error(err)
	} else {
		t.Log(ret)
	}
	if ret, err := getEndpoints("key", "value"); err != nil {
		t.Error(err)
	} else {
		t.Log(ret)
	}
}
