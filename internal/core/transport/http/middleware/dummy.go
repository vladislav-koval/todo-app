package core_http_middleware

import (
	"fmt"
	"net/http"

	core_logger "github.com/vladislav-koval/todo-app/internal/core/logger"
)

func Dummy(s string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())

			log.Debug(fmt.Sprintf("-> before dummy %s", s))

			next.ServeHTTP(w, r)

			log.Debug(fmt.Sprintf("<- after dummy %s", s))
		})
	}
}
