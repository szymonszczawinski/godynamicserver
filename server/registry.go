package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

type serviceRegistry struct {
	ds       *DServer
	outgoing chan string
	services []service
}

func NewServiceRegistry(ds *DServer) *serviceRegistry {
	return &serviceRegistry{
		services: []service{},
		ds:       ds,
		outgoing: make(chan string, 10),
	}
}

func (sr serviceRegistry) GetPort() int {
	return 8080
}

func (sr serviceRegistry) DoGet(w http.ResponseWriter, r *http.Request) error {
	requestContext := NewRequestContext(r)
	slog.Info("registry GET", "path", requestContext.GetPath())
	// requestContext.SetResponseCode(200)
	services := []map[string]any{}
	for _, s := range sr.services {
		services = append(services, s.ToMap())
	}
	responseBody := map[string]any{"services": services}

	// requestContext.SetResponseBody(responseBody)
	sendJSONMessage(w, 200, responseBody)
	return nil
}

func (sr *serviceRegistry) DoPost(w http.ResponseWriter, r *http.Request) error {
	requestContext := NewRequestContext(r)
	slog.Info("registry POST", "path", requestContext.GetPath(), "body", requestContext.GetRequestBody())
	bodyMap := requestContext.GetRequestBody()
	id, ok := bodyMap["id"].(string)
	if !ok {
		responseBody := map[string]any{"error": fmt.Sprintf("id is not a string %v", bodyMap["id"])}
		// requestContext.SetResponseCode(400)
		// requestContext.SetResponseBody(responseBody)
		sendJSONMessage(w, 400, responseBody)
		return nil

	}
	name, ok := bodyMap["name"].(string)
	if !ok {
		responseBody := map[string]any{"error": fmt.Sprintf("name is not a string %v", bodyMap["name"])}
		// requestContext.SetResponseCode(400)
		// requestContext.SetResponseBody(responseBody)
		sendJSONMessage(w, 400, responseBody)
		return nil
	}
	if requestContext.GetPath() == "/" {
		port, ok := bodyMap["port"].(float64)
		if !ok {
			responseBody := map[string]any{"error": fmt.Sprintf("port is not a number %v %T", bodyMap["port"], bodyMap["port"])}
			// requestContext.SetResponseCode(400)
			// requestContext.SetResponseBody(responseBody)
			sendJSONMessage(w, 400, responseBody)
			return nil
		}
		service := NewService(id, name, int(port))
		err := sr.registerService(service)
		if err != nil {
			responseBody := map[string]any{"error": err.Error()}
			// requestContext.SetResponseCode(400)
			// requestContext.SetResponseBody(responseBody)
			sendJSONMessage(w, 400, responseBody)
			return nil
		} else {
			responseBody := map[string]any{"status": "OK"}

			// requestContext.SetResponseCode(201)
			// requestContext.SetResponseBody(responseBody)
			sendJSONMessage(w, 400, responseBody)
			return nil
		}

	} else {
		slog.Warn("registry post path not handled", "path", requestContext.GetPath())
		responseBody := map[string]any{"error": "request path not handled"}
		// requestContext.SetResponseCode(400)
		// requestContext.SetResponseBody(responseBody)
		sendJSONMessage(w, 400, responseBody)
		return nil
	}
}

func (sr *serviceRegistry) OnWebSocketMessage(message string) error {
	slog.Debug("service registry -> on web socket message", "msg", message)
	return nil
}

func (sr *serviceRegistry) Notify(message string) {
	slog.Debug("service registry -> notify", "msg", message)
	sr.outgoing <- message
}

func (sr serviceRegistry) GetOutgoingMessagesQueue() <-chan string {
	return sr.outgoing
}

func (sr *serviceRegistry) registerService(service *service) error {
	slog.Info("registry register service", "service", service)
	sr.services = append(sr.services, *service)
	slog.Info("registry service registerred", "services", sr.services)
	sr.ds.addConnector(service)
	return nil
}
