package server

import (
	"fmt"
	"godynamicserver/model"
	"godynamicserver/service"
	"log/slog"
)

type ServiceRegistry struct {
	ds       *DServer
	services []model.Service
}

func NewServiceRegistry(ds *DServer) *ServiceRegistry {
	return &ServiceRegistry{
		services: []model.Service{},
		ds:       ds,
	}
}

func (sr ServiceRegistry) GetPort() int {
	return 8080
}

func (sr ServiceRegistry) DoGet(requestContext *service.RequestContext) {
	slog.Info("registry  get", "path", requestContext.GetPath())
	requestContext.SetResponseCode(200)
	requestContext.SetResponseBody(sr.services)
	slog.Info("registry get", "services", sr.services)
}

func (sr *ServiceRegistry) DoPost(requestContext *service.RequestContext) {
	slog.Info("registry post", "path", requestContext.GetPath(), "body", requestContext.GetRequestBody())
	bodyMap := requestContext.GetRequestBody()
	id, ok := bodyMap["id"].(string)
	if !ok {
		requestContext.SetResponseCode(400)
		requestContext.SetResponseBody(fmt.Sprintf("id is not a string %v", bodyMap["id"]))

	}
	name, ok := bodyMap["name"].(string)
	if !ok {
		requestContext.SetResponseCode(400)
		requestContext.SetResponseBody(fmt.Sprintf("name is not a string %v", bodyMap["name"]))
	}
	if requestContext.GetPath() == "/" {
		port, ok := bodyMap["port"].(float64)
		if !ok {
			requestContext.SetResponseCode(400)
			requestContext.SetResponseBody(fmt.Sprintf("port is not a number %v %T", bodyMap["port"], bodyMap["port"]))
		}
		service := model.NewService(id, name, int(port))
		err := sr.registerService(service)
		if err != nil {
			requestContext.SetResponseCode(400)
			requestContext.SetResponseBody(err.Error())
		} else {

			requestContext.SetResponseCode(201)
			requestContext.SetResponseBody("OK")
		}

	} else {
		slog.Warn("registry post path not handled", "path", requestContext.GetPath())
		requestContext.SetResponseCode(400)
		requestContext.SetResponseBody("request path not handled")
	}
}

func (sr *ServiceRegistry) registerService(service *model.Service) error {
	slog.Info("registry register service", "service", service)
	sr.services = append(sr.services, *service)
	slog.Info("registry service registerred", "services", sr.services)
	sr.ds.addConnector(service)
	return nil
}
