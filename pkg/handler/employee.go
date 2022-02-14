// Package handler contains application handlers
package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/ivorscott/employee-service/pkg/model"
	"github.com/ivorscott/employee-service/pkg/repository"
	"github.com/ivorscott/employee-service/pkg/trace"
	"github.com/ivorscott/employee-service/pkg/web"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type employeeService interface {
	GetEmployeeByID(ctx context.Context, id string) (model.Employee, error)
	UpdateEmployee(ctx context.Context, employee model.UpdateEmployee) (model.Employee, error)
}

// EmployeeHandler provides method handlers for employee selection.
type EmployeeHandler struct {
	logger    *zap.Logger
	service   employeeService
	publisher rabbitmqAdapter
}

// NewEmployeeHandler creates a new employee handler.
func NewEmployeeHandler(logger *zap.Logger, service employeeService, publisher rabbitmqAdapter) *EmployeeHandler {
	return &EmployeeHandler{
		logger:    logger,
		service:   service,
		publisher: publisher,
	}
}

// GetEmployee retrieves an employee.
func (eh *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.NewSpan(r.Context(), "handler.employee.GetEmployee", nil)
	defer span.End()

	vars := mux.Vars(r)
	test, ok := vars["employee_id"]
	if !ok {
		return web.NewRequestError(repository.ErrInvalidID, http.StatusBadRequest)
	}
	e, err := eh.service.GetEmployeeByID(ctx, test)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidID):
			trace.AddSpanError(span, err)
			return web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, repository.ErrNotFound):
			trace.AddSpanError(span, err)
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			trace.AddSpanError(span, err)
			return err
		}
	}

	return web.Respond(ctx, w, e, http.StatusOK)
}

// UpdateEmployee updates an employee.
func (eh *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.NewSpan(r.Context(), "handler.employee.UpdateSalary", nil)
	defer span.End()

	// call service layer and handle result

	// publish message
	// e.g., publisher.Publish(...)

	return web.Respond(ctx, w, nil, http.StatusOK)
}
