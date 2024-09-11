package server

import (
	"errors"
	"log"

	"golang.org/x/sync/errgroup"
)

var (
	Error404NotFound            = errors.New("not found")
	Error400BadRequest          = errors.New("bad request")
	ErrorUnsupportedMessageType = errors.New("unsupported message type")
)

type DServer struct {
	connectors map[int]*serverConnector
	g          errgroup.Group
}

func NewDServer() *DServer {
	return &DServer{
		connectors: map[int]*serverConnector{},
	}
}

func (ds *DServer) Start() {
	serviceRegistry := NewServiceRegistry(ds)
	ds.addConnector(serviceRegistry)
	if err := ds.g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func (ds *DServer) addConnector(service IService) {
	connector := NewServerConnector(service)
	if _, exist := ds.connectors[service.GetPort()]; !exist {
		ds.connectors[service.GetPort()] = connector
		ds.g.Go(func() error {
			return connector.start()
		})
	}
}
