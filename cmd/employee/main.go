package main

import (
	"context"
	"employee-service/pkg/handlers"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	logger := log.New(os.Stdout, "employee-service: ", log.Lmsgprefix|log.Lmicroseconds|log.Lshortfile)

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
		logger.Println("Starting server...")
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error on startup : %w", err)
	case sig := <-shutdown:
		logger.Printf("Start shutdown due to %q signal\n", sig)

		// give on going tasks a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			fmt.Errorf("server error on startup : %w", err)
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
