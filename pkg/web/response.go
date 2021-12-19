package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// Respond send a response back to the client.
func Respond(ctx context.Context, w http.ResponseWriter, val interface{}, statusCode int) error {
	if v, ok := ctx.Value(KeyValues).(*Values); ok {
		v.StatusCode = statusCode
	}

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}
	res, err := json.Marshal(val)
	if err != nil {
		return err
	}

	w.WriteHeader(statusCode)

	if _, err := w.Write(res); err != nil {
		return err
	}

	return nil
}

// RespondError sends an error response back to the client.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {
	var webErr *Error
	if ok := errors.Is(err, webErr); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}

		return Respond(ctx, w, er, webErr.Status)
	}
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	return Respond(ctx, w, er, http.StatusInternalServerError)
}
