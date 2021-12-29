// Package db maintains the database connection and extensions.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/ivorscott/employee-service/pkg/config"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"github.com/pkg/errors"
)

// Repository represents the database and query builder methods.
type Repository struct {
	SqlxStorer
	SquirrelBuilder
	URL url.URL
	DB  *sql.DB
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
		DB:              db.DB,
		SqlxStorer:      db,
		SquirrelBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
		URL:             u,
	}

	return r, db.Close, nil
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db Storer) error {
	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowxContext(ctx, q).Scan(&tmp)
}

// Storer represents a repository.
type Storer interface {
	SqlxStorer
	SquirrelBuilder
}

// SquirrelBuilder represents the fluent sql generation query builder.
type SquirrelBuilder interface {
	Select(columns ...string) squirrel.SelectBuilder
	Insert(into string) squirrel.InsertBuilder
	Replace(into string) squirrel.InsertBuilder
	Update(table string) squirrel.UpdateBuilder
	Delete(from string) squirrel.DeleteBuilder
	PlaceholderFormat(f squirrel.PlaceholderFormat) squirrel.StatementBuilderType
	RunWith(runner squirrel.BaseRunner) squirrel.StatementBuilderType
}

// SqlxStorer represents the database extension sqlx.
type SqlxStorer interface {
	DriverName() string
	MapperFunc(mf func(string) string)
	Rebind(query string) string
	Unsafe() *sqlx.DB
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	MustBegin() *sqlx.Tx
	Beginx() (*sqlx.Tx, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	MustExec(query string, args ...interface{}) sql.Result
	Preparex(query string) (*sqlx.Stmt, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	MustBeginTx(ctx context.Context, opts *sql.TxOptions) *sqlx.Tx
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Close() error
}
