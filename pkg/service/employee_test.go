package service_test

import (
	"testing"

	"github.com/ivorscott/employee-service/pkg/mocks"
	"github.com/ivorscott/employee-service/pkg/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestEmployeeService_GetEmployeeByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, deps := setupEmployeeServiceTest()
		expected := testEmployees[0]

		deps.repo.On("FindEmployeeByID", mock.AnythingOfType("*context.valueCtx"), expected.ID).Return(expected, nil)
		actual, err := svc.GetEmployeeByID(testCtx, expected.ID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		deps.repo.AssertNumberOfCalls(t, "FindEmployeeByID", 1)
	})
}

type employeeServiceDeps struct {
	logger *zap.Logger
	repo   *mocks.EmployeeRepository
}

func setupEmployeeServiceTest() (*service.EmployeeService, employeeServiceDeps) {
	logger := zap.NewNop()
	repo := &mocks.EmployeeRepository{}
	return service.NewEmployeeService(logger, repo), employeeServiceDeps{logger, repo}
}
