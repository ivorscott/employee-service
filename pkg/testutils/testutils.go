// Package testutils contains test helpers.
package testutils

import (
	"database/sql"

	"github.com/go-testfixtures/testfixtures/v3"
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
