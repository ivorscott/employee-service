// Package testutils contains test helpers.
package testutils

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required by testfixtures

	"github.com/ivorscott/employee-service/pkg/config"
	"github.com/ivorscott/employee-service/pkg/db"
	"github.com/ivorscott/employee-service/res"
)

const (
	dbDriver    = "postgres"
	fixturesDir = "../../res/fixtures"
)

// DBConnect creates a test database connection, migrates up and loads fixtures.
func DBConnect() (*db.Repository, func() error) {
	cfg, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}

	cfg.DB.Name = cfg.TestDB.Name
	cfg.DB.Port = cfg.TestDB.Port
	cfg.DB.DisableTLS = cfg.TestDB.DisableTLS

	repo, Close, err := db.NewRepository(cfg)
	if err != nil {
		panic(err)
	}

	err = res.MigrateUp(repo.URL.String())
	if err != nil {
		panic(err)
	}

	err = loadFixtures(repo.DB)
	if err != nil {
		panic(err)
	}

	return repo, Close
}

func loadFixtures(db *sql.DB) error {
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect(dbDriver),
		testfixtures.Directory(fixturesDir),
	)
	if err != nil {
		return err
	}

	return fixtures.Load()
}

// RunInTx runs test function in a database transaction.
func RunInTx(fn func(t *testing.T, tx *sqlx.Tx)) func(t *testing.T) {
	return func(t *testing.T) {
		conn, Close := DBConnect()
		defer Close()

		conn.RunInTransaction(context.Background(), func(tx *sqlx.Tx) error {
			fn(t, tx)
			return nil
		})
	}
}
