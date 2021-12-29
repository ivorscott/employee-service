// Package service contains all application services.
package service

import (
	"context"

	"github.com/ivorscott/employee-service/pkg/model"
)

type employeeRepository interface {
	FindEmployeeByID(ctx context.Context, id string) (model.Employee, error)
}
