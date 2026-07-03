package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_http_request "github.com/vladislav-koval/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vladislav-koval/todo-app/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

func (h *TaskHttpHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	taskId, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskId path value")
		return
	}

	task, err := h.taskService.GetTask(r.Context(), taskId)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task")

		return
	}

	res := GetTaskResponse(taskDTOFromDomain(task))

	responseHandler.JSONResponse(res, http.StatusOK)
}
