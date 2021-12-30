//go:generate mockery -all -dir . -case snake -output ../mocks
package handler_test

import (
	"os"
	"testing"

	"github.com/ivorscott/employee-service/pkg/model"
	"github.com/ivorscott/employee-service/pkg/testutils"
)

var (
	testEmployees []model.Employee
)

func TestMain(m *testing.M) {
	testutils.LoadGoldenFile(&testEmployees, "employee.json")
	os.Exit(m.Run())
}
