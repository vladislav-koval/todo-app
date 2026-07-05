package web_transport_http

import (
	"net/http"

	core_http_server "github.com/vladislav-koval/todo-app/internal/core/transport/http/server"
)

type WebHttpHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() ([]byte, error)
}

func NewWebHttpHandler(webService WebService) *WebHttpHandler {
	return &WebHttpHandler{
		webService: webService,
	}
}

func (h *WebHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/{$}",
			Handler: h.GetMainPage,
		},
	}
}
