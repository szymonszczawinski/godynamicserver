package model

type Service struct {
	resources map[string]Resource
	id        string
	name      string
	port      int
}

func NewService(id, name string, port int) Service {
	return Service{
		id:   id,
		name: name,
		port: port,
	}
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
