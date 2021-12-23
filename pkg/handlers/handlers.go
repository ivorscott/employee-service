// Package handlers contains application handlers and routes.
package handlers

import (
	"github.com/ivorscott/employee-service/pkg/middleware"
	"github.com/ivorscott/employee-service/pkg/web"
	"go.uber.org/zap"

	"net/http"
	"os"
)

// API configures the application routes, middleware and handlers.
func API(shutdown chan os.Signal, log *zap.Logger) http.Handler {
	e := Employee{}

	app := web.NewApp(shutdown, log, middleware.Logger(log), middleware.Errors(), middleware.Panics(log))
	app.Handle("GET", "/employees/{employee-id}", e.GetEmployee)

	return app
}
