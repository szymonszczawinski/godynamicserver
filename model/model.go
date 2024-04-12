package model

type ServiceDefinition struct {
	resources map[string]ResourceDefinition
	id        string
	name      string
	uri       string
}

func NewServiceDefinition(id, name, uri string) ServiceDefinition {
	return ServiceDefinition{
		id:   id,
		name: name,
		uri:  uri,
	}
}

func (sd *ServiceDefinition) AddResource(id, name string) {
	sd.resources[name] = NewResourceDefiniton(id, name)
}

type ResourceDefinition struct {
	id   string
	name string
}

func NewResourceDefiniton(id, name string) ResourceDefinition {
	return ResourceDefinition{
		id:   id,
		name: name,
	}
}
