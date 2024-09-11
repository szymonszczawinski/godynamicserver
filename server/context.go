package server

import "net/http"

type requestContext struct {
	request      *http.Request
	responseBody map[string]any
	parameters   map[string]any
	requestBody  map[string]any
	path         string
	responseCode int
}

func NewRequestContext(r *http.Request) requestContext {
	return requestContext{
		request: r,
		path:    r.URL.Path,
	}
}

func (rc *requestContext) SetRequestBody(requestBody map[string]any) {
	rc.requestBody = requestBody
}

func (rc requestContext) GetRequestBody() map[string]any {
	return rc.requestBody
}

func (rc requestContext) GetResponseCode() int {
	return rc.responseCode
}

func (rc *requestContext) SetResponseCode(responseCode int) {
	rc.responseCode = responseCode
}

func (rc *requestContext) SetResponseBody(responseBody map[string]any) {
	rc.responseBody = responseBody
}

func (rc requestContext) GetPath() string {
	return rc.path
}

func (rc requestContext) GetResponseBody() map[string]any {
	return rc.responseBody
}
