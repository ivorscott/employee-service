package main

import (
	"embed"
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
func API(
	shutdown chan os.Signal,
	logger *zap.Logger,
	content embed.FS,
	employeeHandler *handler.EmployeeHandler,
) http.Handler {
	mid := []web.Middleware{
		middleware.Metric(),
		middleware.Logger(logger),
		middleware.Error(),
		middleware.Panic(logger),
	}

	swagger, _ := fs.Sub(content, "static/swagger-ui")

	router := mux.NewRouter()
	router.Path("/metrics").Handler(promhttp.Handler())
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.FS(swagger))))

	app := web.NewApp(router, shutdown, logger, mid...)
	app.Handle("GET", "/employees/{employee_id}", "find employee", employeeHandler.GetEmployee)
	app.Handle("PATCH", "/employees/{employee_id}", "update employee", employeeHandler.UpdateEmployee)

	return app
}
