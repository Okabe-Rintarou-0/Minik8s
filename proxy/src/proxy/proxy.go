package proxy

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/proxy/src/constant"
	"minik8s/proxy/src/service"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"strconv"
)

var log = logger.Log("Proxy")

type Proxy struct {
	iptablesManager service.Manager

	endpointUpdates chan *entity.EndpointUpdate
	serviceUpdates  chan *entity.ServiceUpdate
}

func New() (*Proxy, error) {
	if im, err := service.New(); err != nil {
		return nil, err
	} else {
		if err := im.Init(); err != nil {
			return nil, err
		}
		return &Proxy{
			iptablesManager: im,
			serviceUpdates:  make(chan *entity.ServiceUpdate, 20),
			endpointUpdates: make(chan *entity.EndpointUpdate, 20),
		}, nil
	}
}

func (proxy *Proxy) parseEndpointUpdate(msg *redis.Message) {
	endpointUpdate := &entity.EndpointUpdate{}
	err := json.Unmarshal([]byte(msg.Payload), endpointUpdate)
	if err != nil {
		log(err.Error())
		return
	}
	log("Received endpoint update action: %+v", endpointUpdate)
	proxy.endpointUpdates <- endpointUpdate
}

func (proxy *Proxy) parseServiceUpdate(msg *redis.Message) {
	serviceUpdate := &entity.ServiceUpdate{}
	err := json.Unmarshal([]byte(msg.Payload), serviceUpdate)
	if err != nil {
		log(err.Error())
		return
	}
	log("Received service update action: %+v", serviceUpdate)
	proxy.serviceUpdates <- serviceUpdate
}

func hash(s string) string {
	ret := ""
	var i int
	for i = 0; i < 8 && i < len(s); i = i + 1 {
		ret = ret + string(s[i])
	}
	for ; i < len(s); i = i + len(s)/5 {
		ret = ret + string(s[i])
	}
	return ret
}

func hashWithPort(s, port string) string {
	return hash(s) + "-" + port
}

func transEndpoint(endpoint apiObject.Endpoint) []service.EndPoint {
	eps := make([]service.EndPoint, 1)
	eps[0] = service.EndPoint{
		Name: constant.KubeEndpoint + hash(endpoint.UID),
		Ip:   endpoint.IP,
		Port: endpoint.Port,
	}
	return eps
}

func transEndpoints(endpoints []apiObject.Endpoint, expect string) []service.EndPoint {
	eps := make([]service.EndPoint, 0)
	for _, endpoint := range endpoints {
		if endpoint.Port != expect {
			continue
		}
		eps = append(eps, service.EndPoint{
			Name: constant.KubeEndpoint + hash(endpoint.UID),
			Ip:   endpoint.IP,
			Port: endpoint.Port,
		})
	}
	return eps
}

func (proxy *Proxy) syncLoopIteration(endpointUpdates <-chan *entity.EndpointUpdate,
	serviceUpdates <-chan *entity.ServiceUpdate) bool {
	log("Sync loop Iteration")
	select {
	case endpointUpdate := <-endpointUpdates:
		log("Received endpointUpdate %+v", endpointUpdate)

		svc := endpointUpdate.Target.Service
		for _, port := range svc.Spec.Ports {
			serviceName := constant.KubeService + hashWithPort(svc.Metadata.UID, strconv.Itoa(int(port.Port)))
			preEps := transEndpoints(endpointUpdate.Target.PreEndpoints, strconv.Itoa(int(port.TargetPort)))
			newEps := transEndpoints(endpointUpdate.Target.NewEndpoints, strconv.Itoa(int(port.TargetPort)))

			switch endpointUpdate.Action {
			case entity.CreateAction:
				if err := proxy.iptablesManager.DeleteEndPoints(serviceName, preEps); err != nil {
					log(err.Error())
				}
				if err := proxy.iptablesManager.CreateEndpoints(serviceName, newEps); err != nil {
					log(err.Error())
				}
			case entity.UpdateAction:
				log("undefined UpdateAction")
			case entity.DeleteAction:
				if err := proxy.iptablesManager.DeleteEndPoints(serviceName, preEps); err != nil {
					log(err.Error())
				}
				if err := proxy.iptablesManager.CreateEndpoints(serviceName, newEps); err != nil {
					log(err.Error())
				}
			}
		}

	case serviceUpdate := <-serviceUpdates:
		log("Received serviceUpdate %+v", serviceUpdate)

		svc := serviceUpdate.Target.Service
		for _, port := range svc.Spec.Ports {
			serviceName := constant.KubeService + hashWithPort(svc.Metadata.UID, strconv.Itoa(int(port.Port)))
			svcIp := svc.Spec.ClusterIP
			eps := transEndpoints(serviceUpdate.Target.Endpoints, strconv.Itoa(int(port.TargetPort)))

			switch serviceUpdate.Action {
			case entity.CreateAction:
				if err := proxy.iptablesManager.CreateService(serviceName, svcIp+"/32", strconv.Itoa(int(port.Port))); err != nil {
					log(err.Error())
				}
				if err := proxy.iptablesManager.CreateEndpoints(serviceName, eps); err != nil {
					log(err.Error())
				}
			case entity.UpdateAction:
				log("undefined UpdateAction")
			case entity.DeleteAction:
				if err := proxy.iptablesManager.DeleteEndPoints(serviceName, eps); err != nil {
					log(err.Error())
				}
				if err := proxy.iptablesManager.DeleteService(serviceName, svcIp+"/32", strconv.Itoa(int(port.Port))); err != nil {
					log(err.Error())
				}
			}
		}
	}
	return true
}

func (proxy *Proxy) syncLoop(endpointUpdates <-chan *entity.EndpointUpdate,
	serviceUpdates <-chan *entity.ServiceUpdate) {
	for proxy.syncLoopIteration(endpointUpdates, serviceUpdates) {
	}
}

func (proxy *Proxy) Run() {

	go listwatch.Watch(topicutil.EndpointUpdateTopic(), proxy.parseEndpointUpdate)
	go listwatch.Watch(topicutil.ServiceUpdateTopic(), proxy.parseServiceUpdate)

	proxy.syncLoop(proxy.endpointUpdates, proxy.serviceUpdates)
}
