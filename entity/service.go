package entity

import "minik8s/apiObject"

type EndpointTarget struct {
	PreEndpoints []apiObject.Endpoint
	NewEndpoints []apiObject.Endpoint
	Service      apiObject.Service
}

type ServiceTarget struct {
	Endpoints []apiObject.Endpoint
	Service   apiObject.Service
}
