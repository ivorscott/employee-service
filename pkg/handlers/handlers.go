// Package handlers contains application handlers and routes.
package handlers

import (
	"employee-service/pkg/middleware"
	"employee-service/pkg/web"
	"log"
	"net/http"
	"os"
)

// API configures the application routes, middleware and handlers.
func API(shutdown chan os.Signal, log *log.Logger) http.Handler {
	e := Employee{}

	app := web.NewApp(shutdown, log, middleware.Logger(log), middleware.Errors(log), middleware.Panics(log))
	app.Handle("GET", "/employees/{employee-id}", e.GetEmployee)

	return app
}
