package repository

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/ivorscott/employee-service/pkg/model"
	"github.com/ivorscott/employee-service/res/database"
)

// EmployeeRepository is responsible for storing and retrieving employees.
type EmployeeRepository struct {
	repo *database.Repository
}

// NewEmployeeRepository creates a new EmployeeRepository.
func NewEmployeeRepository(repo *database.Repository) *EmployeeRepository {
	return &EmployeeRepository{
		repo: repo,
	}
}

// FindEmployeeByID finds an employee record by id.
func (er *EmployeeRepository) FindEmployeeByID(ctx context.Context, id string) (model.Employee, error) {
	var e model.Employee

	stmt := er.repo.Select(
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

	if err := er.repo.Get(&e, query, id); err != nil {
		if err == sql.ErrNoRows {
			return e, ErrNotFound
		}
		return e, err
	}

	return e, nil
}
