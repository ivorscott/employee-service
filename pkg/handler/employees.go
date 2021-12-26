// Package handler contains application handlers
package handler

import (
	"github.com/ivorscott/employee-service/pkg/model"
	"github.com/ivorscott/employee-service/pkg/web"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Employee service.
type Employee struct{}

// GetEmployee retrieves an employee.
func (e Employee) GetEmployee(w http.ResponseWriter, r *http.Request) error {
	var errMsg struct {
		Error string `json:"error"`
	}

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

	id, err := strconv.Atoi(vars["employee-id"])
	if err != nil {
		errMsg.Error = "bad request"
		return web.Respond(r.Context(), w, errMsg, http.StatusBadRequest)
	}

	if employee, ok := eMap[id]; ok {
		return web.Respond(r.Context(), w, employee, http.StatusOK)
	}
	errMsg.Error = "not found"
	return web.Respond(r.Context(), w, errMsg, http.StatusNotFound)
}
