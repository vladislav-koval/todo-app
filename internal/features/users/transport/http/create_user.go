package users_transport_http

import (
	"net/http"

	"github.com/vladislav-koval/todo-app/internal/core/domain"
	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_http_request "github.com/vladislav-koval/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vladislav-koval/todo-app/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"omitempty,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserResponse UserDTOResponse

func (h *UsersHttpHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	var req CreateUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain, err := h.usersService.CreateUser(r.Context(), domainFromDto(req))

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	res := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(res, http.StatusCreated)
}

func domainFromDto(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
