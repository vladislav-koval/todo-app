package user_service

import (
	"context"
	"fmt"

	"github.com/vladislav-koval/todo-app/internal/core/domain"
)

func (u *UsersService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user: %w", err)
	}

	user, err := u.usersRepository.CreateUser(ctx, user)

	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
