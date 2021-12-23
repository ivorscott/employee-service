package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ivorscott/employee-service/pkg/handlers"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

const logPath = "./logs/out.log"

var logger *zap.Logger

func main() {
	_, err := os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}

	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout", logPath}

	logger, err := c.Build()
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/out.log",
		MaxSize:    25000, // megabytes
		MaxBackups: 24,
		MaxAge:     1,    // days
		Compress:   true, // disabled by default
	})

	if err := run(logger); err != nil {
		log.Fatal(err)
	}
}

func run(logger *zap.Logger) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handlers.API(shutdown, logger),
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
		logger.Info(fmt.Sprintf("Start shutdown due to %q signal\n", sig))

		// give on going tasks a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
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
