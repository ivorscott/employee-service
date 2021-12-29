// Package migration contains application migration logic.
package migration

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // required for golang-migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // required for golang-migrate
	_ "github.com/lib/pq"                                      // The database driver in use.
)

var (
	errNoChange = errors.New("no change")
)

var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

// Migrate applies the latest database migration.
func Migrate(dbname string, url string) error {
	src := fmt.Sprintf("file://%s", basePath)
	m, err := migrate.New(src, url)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != errNoChange {
		log.Fatal(err)
	}
	return nil
}
