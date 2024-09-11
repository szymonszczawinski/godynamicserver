package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func sendJSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message, Code: code})
}

func sendJSONMessage(w http.ResponseWriter, code int, data map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(DataResponse{Data: data, Code: code})
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func handle(hf handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := hf(w, r); err != nil {
			slog.Error("http request error", "request", r.URL, "err", err)
			sendJSONError(w, 400, err.Error())
		}
	}
}

type serverConnector struct {
	router        http.Handler
	httpServer    *http.Server
	wsconnections map[*webSocketConnection]bool
}

func NewServerConnector(service IService) *serverConnector {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", service.GetPort()),
		Handler: router,
	}
	connector := &serverConnector{
		httpServer:    httpServer,
		router:        router,
		wsconnections: map[*webSocketConnection]bool{},
	}
	router.Handle("/*", handle(handleAll(service, connector)))
	return connector
}

func handleAll(s IService, sc *serverConnector) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		slog.Info("handle all", "service", s)
		if isWebSocketRequest(r) {
			return connectWebSocket(w, r, s, sc)
		}
		method := r.Method
		switch method {
		case "GET":
			return s.DoGet(w, r)
		case "POST":
			return s.DoPost(w, r)
		}
		return nil
	}
}

func (sc *serverConnector) start() error {
	slog.Info("start connector", "address", sc.httpServer.Addr)
	return sc.httpServer.ListenAndServe()
}
