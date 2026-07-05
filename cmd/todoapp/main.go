package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/vladislav-koval/todo-app/internal/core/config"
	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	"github.com/vladislav-koval/todo-app/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/vladislav-koval/todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/vladislav-koval/todo-app/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/vladislav-koval/todo-app/internal/features/statistics/repository/postgres"
	statistics_service "github.com/vladislav-koval/todo-app/internal/features/statistics/service"
	statistics_transport_http "github.com/vladislav-koval/todo-app/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/vladislav-koval/todo-app/internal/features/tasks/repository/postgres"
	tasks_service "github.com/vladislav-koval/todo-app/internal/features/tasks/service"
	tasks_transport_http "github.com/vladislav-koval/todo-app/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/vladislav-koval/todo-app/internal/features/users/repository/postgres"
	users_service "github.com/vladislav-koval/todo-app/internal/features/users/service"
	users_transport_http "github.com/vladislav-koval/todo-app/internal/features/users/transport/http"
	"go.uber.org/zap"

	_ "github.com/vladislav-koval/todo-app/docs"
)

// @title 		Golang TODO API
// @version 	1.0
// @description Todo Application REST-API scheme
// @host 		http://127.0.0.1:5050
// @BasePath 	/api/v1
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("Failed to create logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	config := core_config.NewConfigMust()
	time.Local = config.TimeZone
	logger.Debug("application time zone", zap.Any("zone", config.TimeZone))

	logger.Debug("initializing pgx pool")
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init pgx pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	userRepository := users_postgres_repository.NewUsersRepository(pool)
	userService := users_service.NewUsersService(userRepository)
	usersTransportHttp := users_transport_http.NewUsersHttpHandler(userService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	taskRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(taskRepository)
	tasksTransportHttp := tasks_transport_http.NewTasksHttpHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHttp := statistics_transport_http.NewStatisticsHttpHandler(statisticsService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)

	apiVersionRouterV1.RegisterRoutes(usersTransportHttp.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHttp.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHttp.Routes()...)

	//apiVersionRouterV2 := core_http_server.NewApiVersionRouter(
	//core_http_server.ApiVersion2,
	//core_http_middleware.Dummy("api v2 middleware"),
	//)
	//apiVersionRouterV2.RegisterRoutes(usersTransportHttp.Routes()...)

	httpServer.RegisterApiRoutes(apiVersionRouterV1)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Failed to start server", zap.Error(err))
	}
}
