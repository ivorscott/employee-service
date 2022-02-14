//go:generate mockery -all -dir . -case snake -output ../mocks
package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/devpies/employee-service/pkg/model"
	"github.com/devpies/employee-service/pkg/testutils"
)

var (
	testCtx       = context.Background()
	testEmployees []model.Employee
)

func TestMain(m *testing.M) {
	testutils.LoadGoldenFile(&testEmployees, "employee.json")
	os.Exit(m.Run())
}
