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
	services := []map[string]any{}
	for _, s := range sr.services {
		services = append(services, s.ToMap())
	}
	responseBody := map[string]any{"code": 200, "services": services}

	requestContext.SetResponseBody(responseBody)
	slog.Info("registry get", "services", sr.services)
}

func (sr *ServiceRegistry) DoPost(requestContext *service.RequestContext) {
	slog.Info("registry post", "path", requestContext.GetPath(), "body", requestContext.GetRequestBody())
	bodyMap := requestContext.GetRequestBody()
	id, ok := bodyMap["id"].(string)
	if !ok {
		responseBody := map[string]any{"code": 400, "error": fmt.Sprintf("id is not a string %v", bodyMap["id"])}
		requestContext.SetResponseCode(400)
		requestContext.SetResponseBody(responseBody)

	}
	name, ok := bodyMap["name"].(string)
	if !ok {
		responseBody := map[string]any{"code": 400, "error": fmt.Sprintf("name is not a string %v", bodyMap["name"])}
		requestContext.SetResponseCode(400)
		requestContext.SetResponseBody(responseBody)
	}
	if requestContext.GetPath() == "/" {
		port, ok := bodyMap["port"].(float64)
		if !ok {
			responseBody := map[string]any{"code": 400, "error": fmt.Sprintf("port is not a number %v %T", bodyMap["port"], bodyMap["port"])}
			requestContext.SetResponseCode(400)
			requestContext.SetResponseBody(responseBody)
		}
		service := model.NewService(id, name, int(port))
		err := sr.registerService(service)
		if err != nil {
			responseBody := map[string]any{"code": 400, "error": err.Error()}
			requestContext.SetResponseCode(400)
			requestContext.SetResponseBody(responseBody)
		} else {
			responseBody := map[string]any{"code": 201, "status": "OK"}

			requestContext.SetResponseCode(201)
			requestContext.SetResponseBody(responseBody)
		}

	} else {
		slog.Warn("registry post path not handled", "path", requestContext.GetPath())
		responseBody := map[string]any{"code": 400, "error": "request path not handled"}
		requestContext.SetResponseCode(400)
		requestContext.SetResponseBody(responseBody)
	}
}

func (sr *ServiceRegistry) registerService(service *model.Service) error {
	slog.Info("registry register service", "service", service)
	sr.services = append(sr.services, *service)
	slog.Info("registry service registerred", "services", sr.services)
	sr.ds.addConnector(service)
	return nil
}
