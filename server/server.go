package server

import (
	"fmt"
	"godynamicserver/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DServer struct {
	connectors []*ServerConnector
}

func NewDServer() *DServer {
	return &DServer{
		connectors: []*ServerConnector{},
	}
}

func (ds *DServer) Start() {
	serviceRegistry := NewServiceRegistry()
	ds.addConnector(serviceRegistry)
}

func (ds *DServer) addConnector(service service.IService) {
	connector := NewServerConnector(service)
	ds.connectors = append(ds.connectors, connector)
	connector.start()
}

type ServerConnector struct {
	router     *gin.Engine
	httpServer *http.Server
}

func NewServerConnector(service service.IService) *ServerConnector {
	router := gin.New()
	router.NoRoute(handleAll(service))
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", service.GetPort()),
		Handler: router,
	}
	return &ServerConnector{
		httpServer: httpServer,
		router:     router,
	}
}

func handleAll(s service.IService) func(c *gin.Context) {
	return func(c *gin.Context) {
		requestContext := service.NewRequestContext(c.Request.URL.Path, map[string]any{})
		method := c.Request.Method
		switch method {
		case "GET":
			s.DoGet(requestContext)
		case "POST":
			bodyAsMap := map[string]any{}
			err := c.ShouldBindJSON(&bodyAsMap)
			if err != nil {
				slog.Error("parse error", "err", err)
			}
			slog.Info("parsed requst", "body", bodyAsMap)
			requestContext.SetRequestBody(bodyAsMap)
			s.DoPost(requestContext)
		}
		c.JSON(requestContext.GetResponseCode(), gin.H{
			"code":    requestContext.GetResponseCode(),
			"message": requestContext.GetResponseBody(),
		})
	}
}

func (sc *ServerConnector) start() {
	sc.httpServer.ListenAndServe()
}