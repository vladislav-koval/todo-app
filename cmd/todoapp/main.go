package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_postgres_pool "github.com/vladislav-koval/todo-app/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/vladislav-koval/todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/vladislav-koval/todo-app/internal/core/transport/http/server"
	user_postgres_repository "github.com/vladislav-koval/todo-app/internal/features/users/repository/postgres"
	user_service "github.com/vladislav-koval/todo-app/internal/features/users/service"
	user_transport_http "github.com/vladislav-koval/todo-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("Failed to create logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing pgx pool")
	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init pgx pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	userRepository := user_postgres_repository.NewUsersRepository(pool)
	userService := user_service.NewUsersService(userRepository)
	usersTransportHttp := user_transport_http.NewUsersHttpHandler(userService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHttp.Routes()...)
	httpServer.RegisterApiRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Failed to start server", zap.Error(err))
	}
}
