package service

type IService interface {
	GetPort() int
	DoGet(requestContext *RequestContext)
	DoPost(requestContext *RequestContext) error
}
