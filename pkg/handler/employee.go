// Package handler contains application handlers
package handler

import (
	"github.com/ivorscott/employee-service/pkg/model"
	"github.com/ivorscott/employee-service/pkg/trace"
	"github.com/ivorscott/employee-service/pkg/web"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Employee service.
type Employee struct{}

// GetEmployee retrieves an employee.
func (e Employee) GetEmployee(w http.ResponseWriter, r *http.Request) error {
	// Create the parent span.
	ctx, span := trace.NewSpan(r.Context(), "handler.GetEmployee", nil)
	defer span.End()

	vars := mux.Vars(r)
	eMap := map[int]model.Employee{
		0: {
			ID:        0,
			FirstName: "Alan",
			LastName:  "Watts",
			Job:       "Philosopher",
		},
		1: {
			ID:        1,
			FirstName: "John",
			LastName:  "Locke",
			Job:       "Philosopher",
		},
	}

	// Some random informative tags.
	trace.AddSpanTags(span, map[string]string{"param": vars["employee-id"]})

	id, err := strconv.Atoi(vars["employee-id"])
	if err != nil {
		trace.AddSpanError(span, err)
		return web.Respond(ctx, w, nil, http.StatusBadRequest)
	}

	if employee, ok := eMap[id]; ok {
		return web.Respond(ctx, w, employee, http.StatusOK)
	}

	trace.AddSpanError(span, err)
	return web.Respond(ctx, w, nil, http.StatusNotFound)
}
