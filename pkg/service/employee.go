package service

import (
	"context"

	"github.com/ivorscott/employee-service/pkg/model"
	"github.com/ivorscott/employee-service/pkg/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// EmployeeService is responsible for managing employee.
type EmployeeService struct {
	logger       *zap.Logger
	employeeRepo employeeRepository
}

// NewEmployeeService creates a new instance of EmployeeService.
func NewEmployeeService(logger *zap.Logger, employeeRepo employeeRepository) *EmployeeService {
	return &EmployeeService{
		logger:       logger,
		employeeRepo: employeeRepo,
	}
}

// GetEmployeeByID find an employee by id.
func (es *EmployeeService) GetEmployeeByID(ctx context.Context, id string) (model.Employee, error) {
	if _, err := uuid.Parse(id); err != nil {
		return model.Employee{}, repository.ErrInvalidID
	}

	return es.employeeRepo.FindEmployeeByID(ctx, id)
}
