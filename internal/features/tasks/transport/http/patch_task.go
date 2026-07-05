package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/vladislav-koval/todo-app/internal/core/domain"
	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_http_request "github.com/vladislav-koval/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vladislav-koval/todo-app/internal/core/transport/http/response"
	core_http_types "github.com/vladislav-koval/todo-app/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"New task title"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"New task description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean"`
}

type PatchTaskResponse TaskDTOResponse

func (p *PatchTaskRequest) Validate() error {
	if !p.Title.Set && !p.Description.Set && !p.Completed.Set {
		return fmt.Errorf("patch must contain at least one field")
	}

	if p.Title.Set {
		if p.Title.Value == nil {
			return fmt.Errorf("`Title` can't be null")
		}

		titleLen := len([]rune(*p.Title.Value))

		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` length must be between 1 and 100")
		}
	}

	if p.Description.Set && p.Description.Value != nil {
		descriptionLen := len([]rune(*p.Description.Value))

		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf("`Description` length must be between 1 and 1000")
		}
	}

	if p.Completed.Set && p.Completed.Value == nil {
		return fmt.Errorf("`Completed` can't be null")
	}

	return nil
}

// PatchTask 	godoc
// @summary 	Изменить задачу
// @description Изменить задачу в системе
// @description ### Логика обновления полей (Three-state logic):
// @description ### 1. **Поле не передано** `description` - поле игнорируется и не меняется
// @description ### 2. **Передано значение** `description: "New description"` - устанавливает новое значение
// @description ### 2. **Передан null** `description: null` - удаляет значение
// @description title и completed не not nullable
// @Tags 		tasks
// @Accept 		json
// @Produce 	json
// @Param 		id path int true "ID изменяемой задачи"
// @Param 		request body PatchTaskRequest true "Тело запроса"
// @Success 	200	{object} PatchTaskResponse "Успешно измененная задача"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 	409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/tasks{id} [patch]
func (h *TaskHttpHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	taskId, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskId path value")
		return
	}

	var req PatchTaskRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	taskPatch := taskPatchFromRequest(req)

	taskDomain, err := h.taskService.PatchTask(r.Context(), taskId, taskPatch)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	res := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(res, http.StatusOK)
}

func taskPatchFromRequest(req PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		req.Title.ToDomain(),
		req.Description.ToDomain(),
		req.Completed.ToDomain(),
	)
}
