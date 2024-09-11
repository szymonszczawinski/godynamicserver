package server

type (
	SubscribeMessage struct {
		Type          string `json:"type"`
		Event         string `json:"event"`
		Autosubscribe bool   `json:"autosubscribe"`
	}
	UnubscribeMessage struct {
		Type          string `json:"type"`
		Event         string `json:"event"`
		Autosubscribe bool   `json:"autosubscribe"`
	}
	DataMessage struct {
		Type  string         `json:"type"`
		Data  map[string]any `json:"data"`
		Event string         `json:"event"`
	}
	ErrorResponse struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}
	DataResponse struct {
		Data map[string]any `json:"data"`
		Code int            `json:"code"`
	}
)
