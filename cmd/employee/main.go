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

	"github.com/devpies/employee-service/pkg/adapter"
	"github.com/devpies/employee-service/pkg/config"
	"github.com/devpies/employee-service/pkg/db"
	"github.com/devpies/employee-service/pkg/handler"
	"github.com/devpies/employee-service/pkg/model"
	"github.com/devpies/employee-service/pkg/msg"
	"github.com/devpies/employee-service/pkg/repository"
	"github.com/devpies/employee-service/pkg/service"
	"github.com/devpies/employee-service/pkg/trace"
	"github.com/devpies/employee-service/res"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
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

	if err = run(logger, cfg); err != nil {
		logger.Panic("", zap.Error(err))
	}
}

func run(logger *zap.Logger, cfg *config.AppConfig) error {
	var err error

	repo, Close, err := db.NewRepository(cfg)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	defer Close()

	if err = res.MigrateUp(repo.URL.String()); err != nil {
		logger.Fatal("", zap.Error(err))
	}

	rabbitConn := fmt.Sprintf("%s:%s@%s", cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host)

	publisher := adapter.NewRabbitMQPublisher(
		logger,
		rabbitConn,
		amqp091.Config{},
		rabbitmq.WithPublisherOptionsLogging,
	)

	employeeRepository := repository.NewEmployeeRepository(repo)
	employeeService := service.NewEmployeeService(logger, employeeRepository)
	employeeHandler := handler.NewEmployeeHandler(logger, employeeService, publisher)

	// Listen to incoming messages
	go func() {
		l := adapter.NewRabbitMQListener(logger, rabbitConn, amqp091.Config{})
		l.Listen(msg.CreateEmployee, "my_queue", []string{"routing_key1"}, func(ctx context.Context, message msg.Message) error {
			e, eerr := msg.UnmarshalCreateEmployeeCommand(message.Data)
			if eerr != nil {
				return eerr
			}
			return employeeService.CreateEmployee(ctx, model.NewEmployee(e.Data))
		})
		l.Listen(msg.EmployeeUpdated, "my_queue", []string{"routing_key1"}, func(ctx context.Context, message msg.Message) error {
			e, eerr := msg.UnmarshalEmployeeUpdatedEvent(message.Data)
			if eerr != nil {
				return eerr
			}
			_, eerr = employeeService.UpdateEmployee(ctx, model.UpdateEmployee(e.Data))
			return eerr
		})
	}()

	// Gracefully shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      API(shutdown, logger, content, employeeHandler),
	}
	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("Starting server...")
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err = <-serverErrors:
		return fmt.Errorf("server error on startup : %w", err)
	case sig := <-shutdown:
		logger.Info(fmt.Sprintf("Start shutdown due to %s signal", sig))

		// give on going tasks a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		err = srv.Shutdown(ctx)
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
