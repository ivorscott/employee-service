package repository_test

import (
	"testing"

	"github.com/ivorscott/employee-service/pkg/repository"
	"github.com/ivorscott/employee-service/pkg/testutils"

	"github.com/stretchr/testify/assert"
)

var (
	testEmployee = "bc4cd1a1-4f0e-4e39-9960-e6b1cfe388db"
)

func TestEmployeeRepository_FindEmployeeByID(t *testing.T) {
	db, Close := testutils.DBConnect()
	defer Close()

	repo := repository.NewEmployeeRepository(db)
	employee, err := repo.FindEmployeeByID(testCtx, testEmployee)

	assert.Nil(t, err)
	assert.Equal(t, testEmployee, employee.ID)
}
