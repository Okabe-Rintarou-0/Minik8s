package service

import (
	"minik8s/apiObject"
	"minik8s/apiserver/src/url"
	"minik8s/nginx"
	"minik8s/util/weaveutil"
	"strconv"
)

type Manager interface {
	StartService(service apiObject.Service) error
	ApplyService(service apiObject.Service, endpoints []apiObject.Endpoint) error
	ShutdownService(service apiObject.Service) error
}

type serviceManager struct {
}

func New() Manager {
	return &serviceManager{}
}

func (sm *serviceManager) StartService(service apiObject.Service) error {
	nm := nginx.New(service.Metadata.UID)
	if err := nm.Start(); err != nil {
		return err
	}
	return weaveutil.WeaveAttach(nm.GetName(), service.Spec.ClusterIP+url.MaskStr)
}

func (sm *serviceManager) ApplyService(service apiObject.Service, endpoints []apiObject.Endpoint) error {
	nm := nginx.New(service.Metadata.UID)
	var servers []nginx.Server
	for _, port := range service.Spec.Ports {
		var locations []nginx.Location
		for _, endpoint := range endpoints {
			if endpoint.Port == strconv.Itoa(int(port.TargetPort)) {
				location := nginx.Location{
					Dest: endpoint.IP + ":" + endpoint.Port,
					Addr: "/",
				}
				locations = append(locations, location)
			}
		}
		server := nginx.Server{
			Locations: locations,
			Port:      int(port.Port),
		}
		servers = append(servers, server)
	}
	if err := nm.ApplyLoadBalance(servers); err != nil {
		return err
	}
	return nm.Reload()
}

func (sm *serviceManager) ShutdownService(service apiObject.Service) error {
	return nginx.New(service.Metadata.UID).Shutdown()
}
