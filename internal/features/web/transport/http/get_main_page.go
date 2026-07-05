package web_transport_http

import (
	"net/http"

	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_http_response "github.com/vladislav-koval/todo-app/internal/core/transport/http/response"
)

func (h *WebHttpHandler) GetMainPage(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	html, err := h.webService.GetMainPage()

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get html file for main page")
		return
	}

	responseHandler.HtmlResponse(html)
}
