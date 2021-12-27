package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//go:embed static
var content embed.FS

func main() {
	logger, Sync := newLoggerOrFail()
	defer Sync()

	if err := run(logger); err != nil {
		logger.Panic("", zap.Error(err))
	}
}

func run(logger *zap.Logger) error {
	cfg, err := newAppConfig()
	if err != nil {
		return err
	}

	repo, rClose, err := newRepository(cfg)
	if err != nil {
		return err
	}
	defer rClose()

	ctx := context.Background()
	_, tClose, err := newTracer()
	if err != nil {
		return err
	}
	defer tClose(ctx)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      API(shutdown, logger, content, repo),
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
