// Package handlers contains application handlers and routes.
package handlers

import (
	"embed"
	"github.com/gorilla/mux"
	"github.com/ivorscott/employee-service/pkg/middleware"
	"github.com/ivorscott/employee-service/pkg/web"
	"go.uber.org/zap"
	"io/fs"
	"net/http"
	"os"
)

// API configures the application routes, middleware and handlers.
func API(shutdown chan os.Signal, log *zap.Logger, content embed.FS) http.Handler {
	e := Employee{}

	router := mux.NewRouter()
	swagger, _ := fs.Sub(content, "static/swagger-ui")

	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.FS(swagger))))

	app := web.NewApp(router, shutdown, log, middleware.Logger(log), middleware.Errors(), middleware.Panics(log))
	app.Handle("GET", "/employees/{employee-id}", e.GetEmployee)

	return app
}
