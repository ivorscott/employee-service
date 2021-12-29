package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ivorscott/employee-service/pkg/config"
	"github.com/ivorscott/employee-service/pkg/db"
	"github.com/ivorscott/employee-service/pkg/repository"
	"github.com/ivorscott/employee-service/pkg/service"
	"github.com/ivorscott/employee-service/pkg/trace"

	"go.uber.org/zap"
)

//go:embed static
var content embed.FS

func main() {
	logger, Sync := newLoggerOrPanic()
	defer Sync()

	cfg, err := config.NewAppConfig()
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}

	repo, Close, err := db.NewRepository(cfg)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	defer Close()

	ctx := context.Background()
	prv, err := trace.NewProvider(trace.ProviderConfig{
		JaegerEndpoint: "http://localhost:14268/api/traces",
		ServiceName:    "employee-service",
		ServiceVersion: "1.0.0",
		Environment:    "dev",
		Disabled:       false,
	})
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	defer prv.Close(ctx)

	if err := run(logger, repo, cfg); err != nil {
		logger.Panic("", zap.Error(err))
	}
}

func run(logger *zap.Logger, repo *db.Repository, cfg *config.AppConfig) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	employeeRepository := repository.NewEmployeeRepository(repo)

	employeeService := service.NewEmployeeService(logger, employeeRepository)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      API(shutdown, logger, content, employeeService),
	}
	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("Starting server...")
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error on startup : %w", err)
	case sig := <-shutdown:
		logger.Info(fmt.Sprintf("Start shutdown due to %s signal", sig))

		// give on going tasks a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			err = srv.Close()
		}

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.New("could not stop server gracefully")
		}
	}

	return nil
}
