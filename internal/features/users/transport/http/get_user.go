package users_transport_http

import (
	"net/http"

	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_http_request "github.com/vladislav-koval/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vladislav-koval/todo-app/internal/core/transport/http/response"
)

type GetUserResponse UserDTOResponse

// GetUser 		godoc
// @summary 	Получить пользователя
// @description Получение существующего пользователя в системе по его ID
// @Tags 		users
// @Produce 	json
// @Param 		id path int true "ID получаемого пользователя"
// @Success 	200 {object} GetUserResponse	"Пользователь успешно найден"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users/{id} [get]
func (h *UsersHttpHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	userId, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId path value")

		return
	}

	userDomain, err := h.usersService.GetUser(r.Context(), userId)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")

		return
	}

	res := GetUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(res, http.StatusOK)
}
