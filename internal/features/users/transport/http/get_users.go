package user_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_http_response "github.com/vladislav-koval/todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/vladislav-koval/todo-app/internal/core/transport/http/utils"
)

type GetUsersRequest struct {
}

type GetUsersResponse []UserDTOResponse

func (h *UsersHttpHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	limit, offset, err := getLimitOffsetQueryParams(r)
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

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_utils.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParam(r, "offset")

	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
