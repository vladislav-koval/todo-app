package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_http_request "github.com/vladislav-koval/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vladislav-koval/todo-app/internal/core/transport/http/response"
)

func (h *TaskHttpHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	taskId, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task id")
		return
	}

	err = h.taskService.DeleteTask(r.Context(), taskId)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete task")
		return
	}

	responseHandler.NoContentResponse()
}
