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
