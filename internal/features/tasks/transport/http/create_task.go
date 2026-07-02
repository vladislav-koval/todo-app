package tasks_transport_http

import (
	"net/http"

	"github.com/vladislav-koval/todo-app/internal/core/domain"
	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_http_request "github.com/vladislav-koval/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vladislav-koval/todo-app/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserID int     `json:"author_user_id" validate:"required"`
}
type CreateTaskResponse TaskDTOResponse

func (h *TaskHttpHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	var req CreateTaskRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	taskDomain, err := h.taskService.CreateTask(r.Context(), domainFromDto(req))

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")
		return
	}

	res := CreateTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(res, http.StatusCreated)
}

func domainFromDto(dto CreateTaskRequest) domain.Task {
	return domain.NewTaskUninitialized(dto.Title, dto.Description, dto.AuthorUserID)
}
