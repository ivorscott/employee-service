// Package handler contains application handlers
package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/devpies/employee-service/pkg/handler"
	"github.com/devpies/employee-service/pkg/mocks"
)

func TestEmployeeHandler_GetEmployee(t *testing.T) {
	basePath := "/employees"

	t.Run("success", func(t *testing.T) {
		expected := testEmployees[0]
		handle, deps := setupEmployeeRouter()

		path := fmt.Sprintf("%s/%s", basePath, expected.ID)

		r := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()

		deps.service.On("GetEmployeeByID", mock.AnythingOfType("*context.valueCtx"), expected.ID).Return(expected, nil)

		handle.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

type employeeHandlerDeps struct {
	logger    *zap.Logger
	service   *mocks.EmployeeService
	publisher *mocks.RabbitmqAdapter
}

func setupEmployeeRouter() (http.Handler, employeeHandlerDeps) {
	router := mux.NewRouter()
	logger := zap.NewNop()
	service := &mocks.EmployeeService{}
	publisher := &mocks.RabbitmqAdapter{}

	employee := handler.NewEmployeeHandler(logger, service, publisher)

	router.Methods("GET").Path("/employees/{employee_id}").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = employee.GetEmployee(w, r)
		})

	return router, employeeHandlerDeps{logger, service, publisher}
}
