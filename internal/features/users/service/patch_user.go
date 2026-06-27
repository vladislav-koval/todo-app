package user_service

import (
	"context"
	"fmt"

	"github.com/vladislav-koval/todo-app/internal/core/domain"
)

func (u *UsersService) PatchUser(ctx context.Context, id int, userPatch domain.UserPatch) (domain.User, error) {
	user, err := u.usersRepository.GetUser(ctx, id)

	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repo: %w", err)
	}

	if err := user.ApplyPatch(userPatch); err != nil {
		return domain.User{}, fmt.Errorf("apply patch: %w", err)
	}

	patchedUser, err := u.usersRepository.PatchUser(ctx, id, user)

	if err != nil {
		return domain.User{}, fmt.Errorf("patch user: %w", err)
	}

	return patchedUser, nil
}
