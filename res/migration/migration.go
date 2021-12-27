package migration

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	ErrNoChange = errors.New("no change")
	//ErrNilVersion     = errors.New("no migration")
	//ErrInvalidVersion = errors.New("version must be >= -1")
	//ErrLocked         = errors.New("database locked")
	//ErrLockTimeout    = errors.New("timeout: can't acquire database lock")
)

var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

func Migrate(dbname string, url string) error {
	src := fmt.Sprintf("file://%s", basePath)
	m, err := migrate.New(src, url)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != ErrNoChange {
		log.Fatal(err)
	}
	return nil
}
