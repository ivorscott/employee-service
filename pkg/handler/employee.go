// Package handler contains application handlers
package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/devpies/employee-service/pkg/model"
	"github.com/devpies/employee-service/pkg/repository"
	"github.com/devpies/employee-service/pkg/trace"
	"github.com/devpies/employee-service/pkg/web"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

type employeeService interface {
	GetEmployeeByID(ctx context.Context, id string) (model.Employee, error)
	UpdateEmployee(ctx context.Context, employee model.UpdateEmployee) (model.Employee, error)
}

// EmployeeHandler provides method handlers for employee selection.
type EmployeeHandler struct {
	logger           *zap.Logger
	service          employeeService
	publisher        rabbitmqAdapter
	verificationAddr string
	verificationHTTP *http.Client
}

// NewEmployeeHandler creates a new employee handler. verificationAddr points
// at the verification-service used to check eligibility on every lookup -
// see cmd/verification for why that hop exists.
func NewEmployeeHandler(logger *zap.Logger, service employeeService, publisher rabbitmqAdapter, verificationAddr string) *EmployeeHandler {
	return &EmployeeHandler{
		logger:           logger,
		service:          service,
		publisher:        publisher,
		verificationAddr: verificationAddr,
		verificationHTTP: &http.Client{
			// otelhttp.NewTransport injects the traceparent header, so this
			// call joins the same trace as the inbound request's span.
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   3 * time.Second,
		},
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

	if err := eh.verifyEmployee(ctx, e.ID); err != nil {
		// Eligibility checks are best-effort for this lab; log and carry on
		// rather than failing the lookup.
		trace.AddSpanError(span, err)
		eh.logger.Warn("verification check failed", zap.Error(err))
	}

	return web.Respond(ctx, w, e, http.StatusOK)
}

// verifyEmployee calls the verification-service to check the employee is
// eligible. It exists so the trace waterfall and Tempo's service graph have a
// second real hop instead of just employee-service and a virtual caller node.
func (eh *EmployeeHandler) verifyEmployee(ctx context.Context, employeeID string) error {
	url := fmt.Sprintf("%s/verify/%s", eh.verificationAddr, employeeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := eh.verificationHTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("verification-service returned status %d", resp.StatusCode)
	}

	return nil
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
