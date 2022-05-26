package proxy

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/proxy/src/service"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
)

var log = logger.Log("Proxy")

type Proxy struct {
	iptablesManager service.Manager

	endpointUpdates chan *entity.EndpointUpdate
	serviceUpdates  chan *entity.ServiceUpdate
}

func New() *Proxy {
	im := service.New()
	return &Proxy{
		iptablesManager: im,
		serviceUpdates:  make(chan *entity.ServiceUpdate, 20),
		endpointUpdates: make(chan *entity.EndpointUpdate, 20),
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

func (proxy *Proxy) syncLoopIteration(endpointUpdates <-chan *entity.EndpointUpdate,
	serviceUpdates <-chan *entity.ServiceUpdate) bool {
	log("Sync loop Iteration")
	select {
	case endpointUpdate := <-endpointUpdates:
		log("Received endpointUpdate %+v", endpointUpdate)

		switch endpointUpdate.Action {
		case entity.CreateAction:
			if err := proxy.iptablesManager.ApplyService(endpointUpdate.Target.Service, endpointUpdate.Target.NewEndpoints); err != nil {
				log(err.Error())
			}
		case entity.UpdateAction:
			log("undefined UpdateAction")
		case entity.DeleteAction:
			if err := proxy.iptablesManager.ApplyService(endpointUpdate.Target.Service, endpointUpdate.Target.NewEndpoints); err != nil {
				log(err.Error())
			}
		}

	case serviceUpdate := <-serviceUpdates:
		log("Received serviceUpdate %+v", serviceUpdate)

		switch serviceUpdate.Action {
		case entity.CreateAction:
			if err := proxy.iptablesManager.StartService(serviceUpdate.Target.Service); err != nil {
				log(err.Error())
			}
			if err := proxy.iptablesManager.ApplyService(serviceUpdate.Target.Service, serviceUpdate.Target.Endpoints); err != nil {
				log(err.Error())
			}
		case entity.UpdateAction:
			log("undefined UpdateAction")
		case entity.DeleteAction:
			if err := proxy.iptablesManager.ShutdownService(serviceUpdate.Target.Service); err != nil {
				log(err.Error())
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
