package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ivorscott/employee-service/pkg/trace"
)

// Respond send a response back to the client.
func Respond(ctx context.Context, w http.ResponseWriter, val interface{}, statusCode int) error {
	_, span := trace.NewSpan(ctx, "web.Respond", nil)
	defer span.End()

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if statusCode >= 400 {
		w.Header().Set("Content-Type", "application/problem+json")
	}

	if v, ok := ctx.Value(KeyValues).(*Values); ok {
		v.StatusCode = statusCode
	}

	if val != nil {
		res, err := json.Marshal(val)
		if err != nil {
			return err
		}

		_, err = w.Write(res)
		if err != nil {
			return err
		}
	}

	return nil
}

// RespondError sends an error response back to the client.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {
	ctx, span := trace.NewSpan(ctx, "web.RespondError", nil)
	defer span.End()

	var webErr *Error

	if ok := errors.Is(err, webErr); ok {
		var event = make(map[string]string)

		for _, v := range webErr.Fields {
			event[v.Field] = v.Error
		}
		trace.AddSpanEvents(span, "errors", event)
		trace.AddSpanError(span, err)

		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}

		return Respond(ctx, w, er, webErr.Status)
	}

	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	trace.AddSpanError(span, err)
	trace.FailSpan(span, "internal error")
	return Respond(ctx, w, er, http.StatusInternalServerError)
}
