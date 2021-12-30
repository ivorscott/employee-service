package repository_test

import (
	"testing"

	"github.com/ivorscott/employee-service/pkg/repository"
	"github.com/ivorscott/employee-service/pkg/testutils"

	"github.com/stretchr/testify/assert"
)

func TestEmployeeRepository_FindEmployeeByID(t *testing.T) {
	db, Close := testutils.DBConnect()
	defer Close()

	t.Run("success", func(t *testing.T) {
		expected := testEmployees[0]
		repo := repository.NewEmployeeRepository(db)

		actual, err := repo.FindEmployeeByID(testCtx, expected.ID)
		assert.Nil(t, err)

		assert.Equal(t, expected, actual)
	})
}
