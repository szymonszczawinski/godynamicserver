package server

import (
	"godynamicserver/model"
	"godynamicserver/service"
	"log/slog"
)

type ServiceRegistry struct {
	services []*model.ServiceDefinition
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: []*model.ServiceDefinition{},
	}
}

func (sr ServiceRegistry) GetPort() int {
	return 8080
}

func (sr ServiceRegistry) DoGet(requestContext *service.RequestContext) {
	slog.Info("registry  get", "path", requestContext.GetPath())
	requestContext.SetResponseCode(201)
	requestContext.SetResponseBody("Hello World")
}

func (sr *ServiceRegistry) DoPost(requestContext *service.RequestContext) {
	slog.Info("registry post", "path", requestContext.GetPath(), "body", requestContext.GetRequestBody())
	if requestContext.GetPath() == "/" {
		sr.registerService()
	}
}

func (sr *ServiceRegistry) registerService() {}
