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
}

// EmployeeHandler provides method handlers for employee selection.
type EmployeeHandler struct {
	logger  *zap.Logger
	service employeeService
}

// NewEmployeeHandler creates a new employee handler.
func NewEmployeeHandler(logger *zap.Logger, service employeeService) *EmployeeHandler {
	return &EmployeeHandler{
		logger:  logger,
		service: service,
	}
}

// GetEmployee retrieves an employee.
func (eh *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.NewSpan(r.Context(), "handler.employee.GetEmployee", nil)
	defer span.End()

	vars := mux.Vars(r)
	e, err := eh.service.GetEmployeeByID(ctx, vars["employee-id"])
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
