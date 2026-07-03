package tasks_service

import (
	"context"
	"fmt"

	"github.com/vladislav-koval/todo-app/internal/core/domain"
	core_pagination "github.com/vladislav-koval/todo-app/internal/core/pagination"
)

func (s *TasksService) GetTasks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error) {
	pagination, err := core_pagination.NewPagination(limit, offset)

	if err != nil {
		return nil, fmt.Errorf("create pagination: %w", err)
	}

	tasks, err := s.tasksRepository.GetTasks(
		ctx,
		userID,
		pagination.Limit,
		pagination.Offset,
	)

	if err != nil {
		return nil, fmt.Errorf("get tasks from repo: %w", err)
	}

	return tasks, nil
}
