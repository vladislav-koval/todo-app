package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vladislav-koval/todo-app/internal/core/domain"
)

func (r *StatisticsRepository) GetStatistics(ctx context.Context, userID *int, from *time.Time, to *time.Time) (domain.Statistics, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
		SELECT
		    COUNT(*) AS tasks_created,
		
		    COUNT(*) FILTER (
		        WHERE completed
		    ) AS tasks_completed,
		
			COUNT(*) FILTER (
			    WHERE completed
			)::float / NULLIF(COUNT(*), 0) * 100 AS tasks_completed_rate,
		
		    AVG(completed_at - created_at) AS tasks_average_completion_time
		FROM todoapp.tasks
	`)

	var args []any
	var conditions []string

	if userID != nil {
		conditions = append(conditions, fmt.Sprintf("author_user_id = $%d", len(args)+1))
		args = append(args, userID)
	}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", len(args)+1))
		args = append(args, from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at < $%d", len(args)+1))
		args = append(args, to)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	rows := r.pool.QueryRow(ctx, queryBuilder.String(), args...)

	var statisticsModel StatisticsModel

	err := rows.Scan(
		&statisticsModel.TasksCreated,
		&statisticsModel.TasksCompleted,
		&statisticsModel.TasksCompletedRate,
		&statisticsModel.TasksAverageCompletionSeconds,
	)

	if err != nil {
		return domain.Statistics{}, fmt.Errorf("scan error: %w", err)
	}

	return statisticsDomainFromModel(statisticsModel), nil
}
