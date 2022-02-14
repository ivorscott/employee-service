// Package web provides a custom web framework.
package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/devpies/employee-service/pkg/trace"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Handler represents a custom http handler that returns an error.
type Handler func(http.ResponseWriter, *http.Request) error

// App represents a new application.
type App struct {
	log      *zap.Logger
	mux      *mux.Router
	mw       []Middleware
	shutdown chan os.Signal
}

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values carries information about each request.
type Values struct {
	StatusCode int
	Start      time.Time
}

// NewApp returns a new app equipped with built-in middleware required for every handler.
func NewApp(router *mux.Router, shutdown chan os.Signal, logger *zap.Logger, middleware ...Middleware) *App {
	return &App{
		log:      logger,
		mux:      router,
		mw:       middleware,
		shutdown: shutdown,
	}
}

// Handle converts our custom handler to the standard library Handler.
func (a *App) Handle(methods string, path string, name string, h Handler) {
	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			Start: time.Now(),
		}

		ctx := r.Context()                          // get original context
		ctx = context.WithValue(ctx, KeyValues, &v) // create a new context with new key/value
		// you can't directly update a request context
		r = r.WithContext(ctx) // create a new request and pass context

		// Catch any propagated error
		if err := h(w, r); err != nil {
			a.log.Error("", zap.Error(fmt.Errorf("error: unhandled error\n %+v", err)))
			if IsShutdown(err) {
				a.SignalShutdown()
			}
		}
	}

	a.mux.NewRoute().Methods(strings.Split(methods, ",")...).Path(path).HandlerFunc(trace.HTTPHandlerFunc(fn, name))
}

// ServeHTTP extends original mux ServeHTTP method.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

// SignalShutdown sends application shutdown signal.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGSTOP
}
