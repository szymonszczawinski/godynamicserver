package server

import (
	"errors"
	"fmt"
	"godynamicserver/service"
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	Error404NotFound   = errors.New("not found")
	Error400BadRequest = errors.New("bad request")
)

type DServer struct {
	connectors map[int]*ServerConnector
	g          errgroup.Group
}

func NewDServer() *DServer {
	return &DServer{
		connectors: map[int]*ServerConnector{},
	}
}

func (ds *DServer) Start() {
	serviceRegistry := NewServiceRegistry(ds)
	ds.addConnector(serviceRegistry)
	if err := ds.g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func (ds *DServer) addConnector(service service.IService) {
	connector := NewServerConnector(service)
	if _, exist := ds.connectors[service.GetPort()]; !exist {
		ds.connectors[service.GetPort()] = connector
		ds.g.Go(func() error {
			return connector.start()
		})
	}
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
			slog.Info("get", "response", requestContext.GetResponseBody())
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

func (sc *ServerConnector) start() error {
	return sc.httpServer.ListenAndServe()
}
