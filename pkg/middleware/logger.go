package middleware

import (
	"employee-service/pkg/web"
	"errors"
	"log"
	"net/http"
	"time"
)

// Logger writes some information about the request to the logs.
func Logger(log *log.Logger) web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {
		// Create the handler that will be attached in the middleware chain.
		h := func(w http.ResponseWriter, r *http.Request) error {
			v, ok := r.Context().Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("web value missing from context")
			}

			err := before(w, r)

			// format: GET /foo
			log.Printf("(%d) : %s %s (%s)",
				v.StatusCode,
				r.Method,
				r.URL.Path,
				time.Since(v.Start),
			)

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return f
}
