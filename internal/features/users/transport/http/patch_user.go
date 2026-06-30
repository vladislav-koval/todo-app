package user_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vladislav-koval/todo-app/internal/core/domain"
	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_http_request "github.com/vladislav-koval/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vladislav-koval/todo-app/internal/core/transport/http/response"
	core_http_types "github.com/vladislav-koval/todo-app/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (p *PatchUserRequest) Validate() error {
	if !p.FullName.Set && !p.PhoneNumber.Set {
		return fmt.Errorf("patch must contain at least one field")
	}

	if p.FullName.Set {
		if p.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be null")
		}

		fullNameLen := len([]rune(*p.FullName.Value))

		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` length must be between 3 and 100")
		}
	}

	if p.PhoneNumber.Set && p.PhoneNumber.Value != nil {
		phoneNumberLen := len([]rune(*p.PhoneNumber.Value))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("`PhoneNumber` length must be between 10 and 15")
		}

		if !strings.HasPrefix(*p.PhoneNumber.Value, "+") {
			return fmt.Errorf("`PhoneNumber` must starts with '+'")
		}

		phone := *p.PhoneNumber.Value

		for i := 1; i < len(phone); i++ {
			char := phone[i]
			if char < '0' || char > '9' {
				return fmt.Errorf("`PhoneNumber` must be a digit")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHttpHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	id, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId path value")
		return
	}

	var req PatchUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request ")
		return
	}

	userPatch := userPatchFromRequest(req)

	userDomain, err := h.usersService.PatchUser(r.Context(), id, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	res := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(res, http.StatusOK)
}

func userPatchFromRequest(r PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		r.FullName.ToDomain(),
		r.PhoneNumber.ToDomain(),
	)
}
