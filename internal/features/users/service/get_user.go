package user_service

import (
	"context"
	"fmt"

	"github.com/vladislav-koval/todo-app/internal/core/domain"
)

func (u *UsersService) GetUser(ctx context.Context, id int) (domain.User, error) {
	user, err := u.usersRepository.GetUser(ctx, id)

	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repo: %w", err)
	}

	return user, nil
}
