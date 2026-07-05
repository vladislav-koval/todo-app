package users_transport_http

import (
	"net/http"

	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_http_request "github.com/vladislav-koval/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vladislav-koval/todo-app/internal/core/transport/http/response"
)

type GetUsersRequest struct {
}

type GetUsersResponse []UserDTOResponse

// GetUsers 	godoc
// @summary 	Список пользователей
// @description Получение списка пользователей с опциональной пагинацией
// @Tags 		users
// @Produce 	json
// @Param 		limit query int false "Размер страницы с пользователями"
// @Param 		offset query int false "Смещение страницы с пользователями"
// @Success 	200 {object} GetUsersResponse	"Список пользователей успешно найден"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users [get]
func (h *UsersHttpHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit/offset query params")
		return
	}

	userDomains, err := h.usersService.GetUsers(r.Context(), limit, offset)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	res := GetUsersResponse(usersDTOFromDomains(userDomains))

	responseHandler.JSONResponse(res, http.StatusOK)
}
