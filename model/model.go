package model

import (
	"godynamicserver/service"
	"log/slog"
)

type Service struct {
	resources map[string]Resource
	id        string
	name      string
	port      int
}

func NewService(id, name string, port int) *Service {
	return &Service{
		id:   id,
		name: name,
		port: port,
	}
}

func (s Service) GetPort() int {
	return s.port
}

func (s Service) DoGet(requestContext *service.RequestContext) {
	slog.Info("get", "service", s.name, "path", requestContext.GetPath())
}

func (s *Service) DoPost(requestContext *service.RequestContext) {
	slog.Info("post", "service", s.name, "path", requestContext.GetPath())
}

func (sd *Service) AddResource(id, name string) {
	sd.resources[name] = NewResource(id, name)
}

type Resource struct {
	id   string
	name string
}

func NewResource(id, name string) Resource {
	return Resource{
		id:   id,
		name: name,
	}
}
