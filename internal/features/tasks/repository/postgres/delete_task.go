package tasks_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/vladislav-koval/todo-app/internal/core/errors"
)

func (r *TasksRepository) DeleteTask(ctx context.Context, taskId int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM todoapp.tasks
		WHERE id = $1;
	`

	ct, err := r.pool.Exec(ctx, query, taskId)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("task with id '%d': %w", taskId, core_errors.ErrNotFound)
	}

	return nil
}
