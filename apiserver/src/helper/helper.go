package helper

import (
	"encoding/json"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"minik8s/util/uidutil"
	"path"
)

func GetNodeHostnames() (hostnames []string) {
	nodeRaws, err := etcd.GetAll(url.NodeURL)
	if err != nil {
		return nil
	}

	for _, nodeRaw := range nodeRaws {
		node := &apiObject.Node{}
		if err = json.Unmarshal([]byte(nodeRaw), node); err == nil {
			hostnames = append(hostnames, node.Name())
		}
	}
	return
}

func ExistsPod(namespace, name string) bool {
	hostnames := GetNodeHostnames()
	var etcdURL string
	for _, hostname := range hostnames {
		etcdURL = path.Join(url.PodURL, hostname, namespace, name)
		if podRaw, err := etcd.Get(etcdURL); err == nil && podRaw != "" {
			return true
		}
	}
	return false
}

func GetPodsApiObjectFromEtcd(node string) (pods []*apiObject.Pod) {
	etcdURL := path.Join(url.PodURL, node)
	if raws, err := etcd.GetAll(etcdURL); err == nil {
		for _, raw := range raws {
			pod := &apiObject.Pod{}
			if err = json.Unmarshal([]byte(raw), &pod); err == nil {
				pods = append(pods, pod)
			}
		}
	}
	return
}

func AddEndpoints(pod apiObject.Pod) error {
	var err error
	for key, value := range pod.Metadata.Labels {
		etcdEndpointsKVURL := path.Join(url.EndpointURL, key, value)
		var endpointsJsonStr string
		if endpointsJsonStr, err = etcd.Get(etcdEndpointsKVURL); err != nil {
			return err
		}
		endpoints := make(etcd.Endpoints)
		if len(endpointsJsonStr) != 0 {
			if err := json.Unmarshal([]byte(endpointsJsonStr), &endpoints); err != nil {
				return err
			}
		}

		preEndpoints := make([]apiObject.Endpoint, 0)
		for _, UIDs := range endpoints {
			for _, UID := range UIDs {
				endpoint := apiObject.Endpoint{}
				if endpointJsonStr, err := etcd.Get(path.Join(url.EndpointURL, UID)); err != nil {
					return nil
				} else {
					if err := json.Unmarshal([]byte(endpointJsonStr), &endpoint); err != nil {
						return nil
					}
				}
				preEndpoints = append(preEndpoints, endpoint)
			}
		}

		newEndpoints := make([]apiObject.Endpoint, 0)
		for _, container := range pod.Spec.Containers {
			for _, port := range container.Ports {
				endpoint := apiObject.Endpoint{
					UID:  uidutil.New(),
					IP:   pod.Spec.ClusterIp,
					Port: port.ContainerPort,
				}
				endpoints.Add(pod.Metadata.UID, endpoint.UID)
				newEndpoints = append(newEndpoints, endpoint)
				var endpointJson []byte
				if endpointJson, err = json.Marshal(endpoint); err != nil {
					return err
				}
				if err = etcd.Put(path.Join(url.EndpointURL, endpoint.UID), string(endpointJson)); err != nil {
					return err
				}
			}
		}
		newEndpoints = append(preEndpoints, newEndpoints...)

		// Push endpoints to proxy
		if servicesStr, err := etcd.GetAll(path.Join(url.ServiceURL, key, value)); err != nil {
			return err
		} else {
			for _, serviceStr := range servicesStr {
				service := apiObject.Service{}
				if err := json.Unmarshal([]byte(serviceStr), &service); err != nil {
					return err
				}
				endpointUpdateMsg, _ := json.Marshal(entity.EndpointUpdate{
					Action: entity.CreateAction,
					Target: entity.EndpointTarget{
						PreEndpoints: preEndpoints,
						NewEndpoints: newEndpoints,
						Service:      service,
					},
				})
				listwatch.Publish(topicutil.EndpointUpdateTopic(), endpointUpdateMsg)
			}
		}

		var endpointsJson []byte
		if endpointsJson, err = json.Marshal(endpoints); err != nil {
			return err
		}
		if err = etcd.Put(etcdEndpointsKVURL, string(endpointsJson)); err != nil {
			return err
		}
	}
	return nil
}

func DelEndpoints(pod apiObject.Pod) error {
	var err error
	for key, value := range pod.Metadata.Labels {
		etcdEndpointsKVURL := path.Join(url.EndpointURL, key, value)
		var endpointsJsonStr string
		if endpointsJsonStr, err = etcd.Get(etcdEndpointsKVURL); err != nil {
			return err
		}
		endpoints := make(etcd.Endpoints)
		if len(endpointsJsonStr) != 0 {
			if err := json.Unmarshal([]byte(endpointsJsonStr), &endpoints); err != nil {
				return err
			}
		}

		preEndpoints := make([]apiObject.Endpoint, 0)
		newEndpoints := make([]apiObject.Endpoint, 0)
		for podUID, UIDs := range endpoints {
			for _, UID := range UIDs {
				endpoint := apiObject.Endpoint{}
				if endpointJsonStr, err := etcd.Get(path.Join(url.EndpointURL, UID)); err != nil {
					return err
				} else {
					if err := json.Unmarshal([]byte(endpointJsonStr), &endpoint); err != nil {
						return err
					}
				}
				preEndpoints = append(preEndpoints, endpoint)
				if podUID != pod.Metadata.UID {
					newEndpoints = append(newEndpoints, endpoint)
				} else {
					// Push endpoints to proxy
					endpoint := apiObject.Endpoint{}
					if endpointJsonStr, err := etcd.Get(path.Join(url.EndpointURL, UID)); err != nil {
						return err
					} else {
						if err := json.Unmarshal([]byte(endpointJsonStr), &endpoint); err != nil {
							return err
						}
					}
					if err = etcd.Delete(path.Join(url.EndpointURL, UID)); err != nil {
						return err
					}
				}
			}
		}

		if servicesStr, err := etcd.GetAll(path.Join(url.ServiceURL, key, value)); err != nil {
			return err
		} else {
			for _, serviceStr := range servicesStr {
				service := apiObject.Service{}
				if err := json.Unmarshal([]byte(serviceStr), &service); err != nil {
					return err
				}
				endpointUpdateMsg, _ := json.Marshal(entity.EndpointUpdate{
					Action: entity.DeleteAction,
					Target: entity.EndpointTarget{
						PreEndpoints: preEndpoints,
						NewEndpoints: newEndpoints,
						Service:      service,
					},
				})
				listwatch.Publish(topicutil.EndpointUpdateTopic(), endpointUpdateMsg)
			}
		}
		endpoints.Del(pod.Metadata.UID)
		var endpointsJson []byte
		if endpointsJson, err = json.Marshal(endpoints); err != nil {
			return err
		}
		if err = etcd.Put(etcdEndpointsKVURL, string(endpointsJson)); err != nil {
			return err
		}
	}
	return nil
}

func GetEndpoints(key, value string) (endpointArray []apiObject.Endpoint, err error) {
	etcdEndpointsKVURL := path.Join(url.EndpointURL, key, value)
	var endpointsJsonStr string
	if endpointsJsonStr, err = etcd.Get(etcdEndpointsKVURL); err != nil {
		return nil, err
	}
	endpoints := make(etcd.Endpoints)
	if len(endpointsJsonStr) != 0 {
		if err := json.Unmarshal([]byte(endpointsJsonStr), &endpoints); err != nil {
			return nil, err
		}
	}

	endpointArray = make([]apiObject.Endpoint, 0)
	for _, arr := range endpoints {
		for _, UID := range arr {
			if endpointStr, err := etcd.Get(path.Join(url.EndpointURL, UID)); err != nil {
				return nil, err
			} else {
				endpoint := apiObject.Endpoint{}
				if err := json.Unmarshal([]byte(endpointStr), &endpoint); err != nil {
					return nil, err
				}
				endpointArray = append(endpointArray, endpoint)
			}
		}
	}
	return endpointArray, nil
}

func ExistsService(namespace, name string) bool {
	var etcdURL string
	etcdURL = path.Join(url.ServiceURL, namespace, name)
	if serviceRaw, err := etcd.Get(etcdURL); err == nil && serviceRaw != "" {
		return true
	}
	return false
}
