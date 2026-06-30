package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	"github.com/vladislav-koval/todo-app/internal/core/repository/postgres/pool/pgx"
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
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
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
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHttp.Routes()...)

	//apiVersionRouterV2 := core_http_server.NewApiVersionRouter(
	//core_http_server.ApiVersion2,
	//core_http_middleware.Dummy("api v2 middleware"),
	//)
	//apiVersionRouterV2.RegisterRoutes(usersTransportHttp.Routes()...)

	httpServer.RegisterApiRoutes(apiVersionRouterV1)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Failed to start server", zap.Error(err))
	}
}
