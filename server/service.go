package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type IService interface {
	GetPort() int
	DoGet(w http.ResponseWriter, r *http.Request) error
	DoPost(w http.ResponseWriter, r *http.Request) error
	OnWebSocketMessage(message string) error
	Notify(message string)
	GetOutgoingMessagesQueue() <-chan string
}
type service struct {
	resources map[string]Resource
	outgoing  chan string
	id        string
	name      string
	port      int
}

func NewService(id, name string, port int) *service {
	return &service{
		id:       id,
		name:     name,
		port:     port,
		outgoing: make(chan string, 10),
	}
}

func (s service) String() string {
	jsonMap := map[string]any{"id": s.id, "name": s.name, "port": s.port}
	jsonString, _ := json.Marshal(jsonMap)
	return string(jsonString)
}

func (s service) ToMap() map[string]any {
	return map[string]any{"id": s.id, "name": s.name, "port": s.port}
}

func (s service) GetPort() int {
	return s.port
}

func (s service) DoGet(w http.ResponseWriter, r *http.Request) error {
	requestContext := NewRequestContext(r)
	slog.Info("get", "service", s.name, "path", requestContext.GetPath())
	return nil
}

func (s *service) DoPost(w http.ResponseWriter, r *http.Request) error {
	requestContext := NewRequestContext(r)
	slog.Info("post", "service", s.name, "path", requestContext.GetPath())
	return nil
}

func (s *service) OnWebSocketMessage(message string) error {
	slog.Debug("on web socket message", "service", s, "msg", message)
	s.Notify(fmt.Sprintf("ACK %v", message))
	return nil
}

func (s *service) Notify(message string) {
	slog.Debug("on web socket message", "service", s, "msg", message)
	s.outgoing <- message
}

func (s service) GetOutgoingMessagesQueue() <-chan string {
	return s.outgoing
}

func (s *service) AddResource(id, name string) {
	s.resources[name] = NewResource(id, name)
}

type Resource struct {
	elements map[string]any
	id       string
	name     string
}

func NewResource(id, name string) Resource {
	return Resource{
		id:       id,
		name:     name,
		elements: map[string]any{},
	}
}
