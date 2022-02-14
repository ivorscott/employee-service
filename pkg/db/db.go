// Package db maintains the database connection and extensions.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/devpies/employee-service/pkg/config"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"github.com/pkg/errors"
)

// Repository represents a database repository.
type Repository struct {
	*sqlx.DB
	SQ  squirrel.StatementBuilderType
	URL url.URL
}

// NewRepository creates a new repository, connecting it to the postgres server.
func NewRepository(cfg *config.AppConfig) (*Repository, func() error, error) {
	// Define SSL mode.
	sslMode := "require"
	if cfg.DB.DisableTLS {
		sslMode = "disable"
	}

	// Query parameters.
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	// Construct url.
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.DB.User, cfg.DB.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.DB.Host, cfg.DB.Port),
		Path:     cfg.DB.Name,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("postgres", u.String())
	if err != nil {
		return nil, nil, errors.Wrap(err, "connecting to database")
	}

	r := &Repository{
		DB:  db,
		SQ:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
		URL: u,
	}

	return r, db.Close, nil
}

// RunInTransaction runs callback function in a transaction.
func (r *Repository) RunInTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := r.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	return txRun(tx, fn)
}

func txRun(tx *sqlx.Tx, fn func(*sqlx.Tx) error) error {
	defer func() {
		if err := recover(); err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("tx.Rollback panicked: %s", err)
			}
			panic(err)
		}
	}()

	if err := fn(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Printf("tx.Rollback failed: %s", err)
		}
		return err
	}
	return tx.Commit()
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *Repository) error {
	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowxContext(ctx, q).Scan(&tmp)
}
