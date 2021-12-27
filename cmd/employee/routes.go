package main

import (
	"embed"
	"github.com/ivorscott/employee-service/res/database"
	"io/fs"
	"net/http"
	"os"

	"github.com/ivorscott/employee-service/pkg/handler"
	"github.com/ivorscott/employee-service/pkg/middleware"
	"github.com/ivorscott/employee-service/pkg/web"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// API configures the application routes, middleware and handlers.
func API(shutdown chan os.Signal, log *zap.Logger, content embed.FS, repo *database.Repository) http.Handler {
	e := handler.Employee{}

	swagger, _ := fs.Sub(content, "static/swagger-ui")

	router := mux.NewRouter()
	router.Path("/metrics").Handler(promhttp.Handler())
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.FS(swagger))))

	app := web.NewApp(router, shutdown, log, middleware.Metric(), middleware.Logger(log), middleware.Error(), middleware.Panic(log))
	app.Handle("GET", "/employees/{employee-id}", "find employee", e.GetEmployee)

	return app
}
