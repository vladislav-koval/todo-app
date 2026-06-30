package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
	core_http_middleware "github.com/vladislav-koval/todo-app/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux         *http.ServeMux
	config      Config
	log         *core_logger.Logger
	middlewares []core_http_middleware.Middleware
}

func NewHTTPServer(config Config, log *core_logger.Logger, middlewares ...core_http_middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:         http.NewServeMux(),
		config:      config,
		log:         log,
		middlewares: middlewares,
	}
}

func (s *HTTPServer) RegisterApiRoutes(routes ...*ApiVersionRouter) {
	for _, router := range routes {
		prefix := fmt.Sprintf("/api/%s", router.apiVersion)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.WithMiddleware()),
		)
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddlewares(s.mux, s.middlewares...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	s.log.Warn("starting http server", zap.String("addr", s.config.Addr))

	go func() {
		defer close(ch)
		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		return fmt.Errorf("http server: %w", err)
	case <-ctx.Done():
		s.log.Warn("shutdown http server", zap.String("addr", s.config.Addr))
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)

		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown http server: %w", err)
		}

		s.log.Warn("http server stopped", zap.String("addr", s.config.Addr))

		return nil
	}
}
