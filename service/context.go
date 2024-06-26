package service

type RequestContext struct {
	responseBody map[string]any
	parameters   map[string]any
	requestBody  map[string]any
	path         string
	responseCode int
}

func NewRequestContext(path string, parameters map[string]any) *RequestContext {
	return &RequestContext{
		path:       path,
		parameters: parameters,
	}
}

func (rc *RequestContext) SetRequestBody(requestBody map[string]any) {
	rc.requestBody = requestBody
}

func (rc RequestContext) GetRequestBody() map[string]any {
	return rc.requestBody
}

func (rc RequestContext) GetResponseCode() int {
	return rc.responseCode
}

func (rc *RequestContext) SetResponseCode(responseCode int) {
	rc.responseCode = responseCode
}

func (rc *RequestContext) SetResponseBody(responseBody map[string]any) {
	rc.responseBody = responseBody
}

func (rc RequestContext) GetPath() string {
	return rc.path
}

func (rc RequestContext) GetResponseBody() map[string]any {
	return rc.responseBody
}
