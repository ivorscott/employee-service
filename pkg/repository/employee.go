package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ivorscott/employee-service/pkg/db"
	"github.com/ivorscott/employee-service/pkg/model"
	"github.com/ivorscott/employee-service/pkg/trace"

	sq "github.com/Masterminds/squirrel"
)

// EmployeeRepository is responsible for storing and retrieving employees.
type EmployeeRepository struct {
	repo *db.Repository
}

// NewEmployeeRepository creates a new EmployeeRepository.
func NewEmployeeRepository(repo *db.Repository) *EmployeeRepository {
	return &EmployeeRepository{
		repo: repo,
	}
}

// FindEmployeeByID finds an employee record by id.
func (er *EmployeeRepository) FindEmployeeByID(ctx context.Context, id string) (model.Employee, error) {
	ctx, span := trace.NewSpan(ctx, "repository.employee.FindEmployeeByID", nil)
	defer span.End()

	var e model.Employee

	stmt := er.repo.SQ.Select(
		"employee_id",
		"auth0_id",
		"email_address",
		"first_name",
		"middle_name",
		"last_name",
		"phone_number",
		"birth_date",
		"hire_date",
		"picture",
		"language",
		"country",
		"city",
		"zipcode",
		"salary",
		"position",
		"emergency_contact_name",
		"emergency_contact_email",
		"emergency_contact_phone",
		"deleted",
		"updated_at",
		"created_at",
	).From(
		"employees",
	).Where(sq.Eq{"employee_id": "?"})

	query, args, err := stmt.ToSql()
	if err != nil {
		return e, fmt.Errorf("%w: arguments (%v)", err, args)
	}

	if err := er.repo.GetContext(ctx, &e, query, id); err != nil {
		if err == sql.ErrNoRows {
			return e, ErrNotFound
		}
		return e, err
	}

	return e, nil
}
