package users_service

import (
	"context"
	"fmt"

	"github.com/vladislav-koval/todo-app/internal/core/domain"
	core_pagination "github.com/vladislav-koval/todo-app/internal/core/pagination"
)

func (s *UsersService) GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error) {
	pagination, err := core_pagination.NewPagination(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("create pagination: %w", err)
	}

	users, err := s.usersRepository.GetUsers(
		ctx,
		pagination.Limit,
		pagination.Offset,
	)

	if err != nil {
		return nil, fmt.Errorf("get users from repo: %w", err)
	}

	return users, nil
}
