package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_http_request "github.com/vladislav-koval/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vladislav-koval/todo-app/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks 	godoc
// @summary 	Список задач
// @description Получение списка задач с опциональной пагинацией и/или фильтрацией по id автора
// @Tags 		tasks
// @Produce 	json
// @Param 		user_id query int false "Фильтрация задач по ID автора"
// @Param 		limit query int false "Размер страницы с задачами"
// @Param 		offset query int false "Смещение страницы с задачами"
// @Success 	200 {object} GetTasksResponse	"Список задач"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/tasks [get]
func (h *TaskHttpHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	userID, limit, offset, err := getUserIdLimitOffsetQueryParam(r)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID/limit/offset query param")
		return
	}

	taskDomains, err := h.taskService.GetTasks(r.Context(), userID, limit, offset)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	res := GetTasksResponse(tasksDTOFromDomains(taskDomains))

	responseHandler.JSONResponse(res, http.StatusOK)
}

func getUserIdLimitOffsetQueryParam(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDParamKey = "user_id"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDParamKey)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)

	if err != nil {
		return nil, nil, nil, err
	}

	return userID, limit, offset, nil
}
