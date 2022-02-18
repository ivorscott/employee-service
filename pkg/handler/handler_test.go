//go:generate mockery --all --dir . --case snake --output ../mocks --exported
package handler_test

import (
	"os"
	"testing"

	"github.com/devpies/employee-service/pkg/model"
	"github.com/devpies/employee-service/pkg/testutils"
)

var (
	testEmployees []model.Employee
)

func TestMain(m *testing.M) {
	testutils.LoadGoldenFile(&testEmployees, "employee.json")
	os.Exit(m.Run())
}
