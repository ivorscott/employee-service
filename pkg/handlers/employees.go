package handlers

import (
	"employee-service/pkg/models"
	"employee-service/pkg/web"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Employee service.
type Employee struct{}

// GetEmployee retrieves an employee.
func (e Employee) GetEmployee(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	eMap := map[int]models.Employee{
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
		return web.Respond(r.Context(), w, nil, http.StatusBadRequest)
	}

	if employee, ok := eMap[id]; ok {
		return web.Respond(r.Context(), w, employee, http.StatusOK)
	}
	return web.Respond(r.Context(), w, nil, http.StatusNotFound)
}
