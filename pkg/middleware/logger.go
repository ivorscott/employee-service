package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/ivorscott/employee-service/pkg/web"
	"go.uber.org/zap"
)

// Logger middleware writes some information about the request to the logs.
func Logger(log *zap.Logger) web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {
		// Create the handler that will be attached in the middleware chain.
		h := func(w http.ResponseWriter, r *http.Request) error {
			v, ok := r.Context().Value(web.KeyValues).(*web.Values)
			if !ok {
				return fmt.Errorf("web value missing from context")
			}

			err := before(w, r)

			fields := []zap.Field{
				zap.Int("status", v.StatusCode),
				zap.String("method", r.Method),
				zap.String("endpoint", r.URL.Path),
				zap.Duration("duration_ms", time.Since(v.Start)),
				zap.String("human_readable_duration", fmt.Sprint(time.Since(v.Start))),
			}

			if v.StatusCode < 500 {
				log.Info("", fields...)
			} else {
				log.Error("", append(fields, zap.ByteString("stack", debug.Stack()))...)
			}

			return err
		}

		return h
	}

	return f
}
