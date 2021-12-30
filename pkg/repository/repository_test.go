package repository_test

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/ivorscott/employee-service/pkg/model"
	"github.com/ivorscott/employee-service/pkg/testutils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	shouldUpdate  *bool
	testCtx       = context.Background()
	testEmployees []model.Employee
)

func TestMain(m *testing.M) {
	shouldUpdate = flag.Bool("update", false, "update golden files")
	flag.Parse()

	testutils.LoadGoldenFile(&testEmployees, "employee.json")

	os.Exit(m.Run())
}

func TestGoldenFiles(t *testing.T) {
	golden := testutils.NewGoldenConfig(shouldUpdate)

	t.Run("employee golden file", testutils.
		RunInTx(func(t *testing.T, tx *sqlx.Tx) {
			var actual []model.Employee
			var expected []model.Employee

			err := tx.Select(&actual, "SELECT * FROM employees")
			require.NoError(t, err)

			goldenFile := "employee.json"

			if golden.ShouldUpdate() {
				testutils.SaveGoldenFile(&actual, goldenFile)
			}

			testutils.LoadGoldenFile(&expected, goldenFile)
			assert.Equal(t, expected, actual)
		}))
}
