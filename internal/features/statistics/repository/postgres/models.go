package statistics_postgres_repository

import (
	"time"

	"github.com/vladislav-koval/todo-app/internal/core/domain"
)

type StatisticsModel struct {
	TasksCreated                  int
	TasksCompleted                int
	TasksCompletedRate            *float64
	TasksAverageCompletionSeconds *time.Duration
}

func statisticsDomainFromModel(statisticsModel StatisticsModel) domain.Statistics {
	return domain.NewStatistics(
		statisticsModel.TasksCreated,
		statisticsModel.TasksCompleted,
		statisticsModel.TasksCompletedRate,
		statisticsModel.TasksAverageCompletionSeconds,
	)
}
