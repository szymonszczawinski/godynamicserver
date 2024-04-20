package server

import (
	"errors"
	"fmt"
	"godynamicserver/model"
	"godynamicserver/service"
	"log/slog"
)

type ServiceRegistry struct {
	services []*model.Service
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: []*model.Service{},
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

func (sr *ServiceRegistry) DoPost(requestContext *service.RequestContext) error {
	slog.Info("registry post", "path", requestContext.GetPath(), "body", requestContext.GetRequestBody())
	bodyMap := requestContext.GetRequestBody()
	id, ok := bodyMap["id"].(string)
	if !ok {
		return fmt.Errorf("id is not a string %v", bodyMap["id"])
	}
	name, ok := bodyMap["name"].(string)
	if !ok {
		return fmt.Errorf("name is not a string %v", bodyMap["name"])
	}
	if requestContext.GetPath() == "/" {
		port, ok := bodyMap["port"].(float64)
		if !ok {
			return fmt.Errorf("port is not a number %v %T", bodyMap["port"], bodyMap["port"])
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
		return errors.Join(Error400BadRequest, fmt.Errorf(" request path not handled %v", requestContext.GetPath()))
	}
	return nil
}

func (sr *ServiceRegistry) registerService(service model.Service) error {
	slog.Info("service registry register service", "service", service)
	return nil
}
