package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/kaphos/webapp/internal/errorhandling"
	"time"
)

const timeout = time.Second * 2

// Query performs a database query, returning a list of rows.
// If only a single row is required, QueryRow should be used instead.
// to standardise with the other 2 functions. Any errors encountered
// internally are automatically handled using the errorhandling package.
func (d *Database) Query(spanName string, parentCtx context.Context, query string, args ...interface{}) (pgx.Rows, func(), error) {
	ctx, span := d.tracer.Start(parentCtx, spanName)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	rows, err := d.pool.Query(ctx, query, args...)
	err = convertUserError(err)
	errorhandling.Check(err, spanName)

	endFn := func() {
		span.End()
		cancel()
	}

	return rows, endFn, err
}

type QueryRowResult struct {
	spanName string
	row      pgx.Row
	end      func()
}

// QueryRow performs a database query and returns a single row.
// Should be preferred over Query if only a single row is needed.
// Should be called directly with Scan. Any errors encountered
// internally are automatically handled using the errorhandling package.
func (d *Database) QueryRow(spanName string, ctx context.Context, query string, args ...interface{}) QueryRowResult {
	ctx, span := d.tracer.Start(ctx, spanName)
	ctx, cancel := context.WithTimeout(ctx, timeout)

	row := d.pool.QueryRow(ctx, query, args...)

	endFn := func() {
		span.End()
		cancel()
	}

	return QueryRowResult{spanName, row, endFn}
}

// Scan the results from QueryRow into the destination interface(s).
// Called as a second function after QueryRow (instead of combining)
// so that "[]interface{}{}" does not need to be passed into the
// function call.
func (r QueryRowResult) Scan(dest ...interface{}) error {
	defer r.end()
	err := convertUserError(r.row.Scan(dest))
	errorhandling.Check(err, r.spanName)
	return err
}

// Exec executes a query against the database, ignoring any results.
// Should be used when a response is not needed (e.g., when performing
// an update or delete operation). The preferred function if a response
// from the database is not required. Any errors encountered
// internally are automatically handled using the errorhandling package.
func (d *Database) Exec(spanName string, ctx context.Context, query string, args ...interface{}) error {
	ctx, span := d.tracer.Start(ctx, spanName)
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := d.pool.Exec(ctx, query, args...)
	err = convertUserError(err)
	errorhandling.Check(err, spanName)
	return err
}

// TODO: Implement functions for transactions
// Will be good if we can create a helper to instantiate the transaction,
// and group everything else together.

//func (d *Database) txExec(tx pgx.Tx, opts Opts) errchk {
//	ctx, cancel := context.WithTimeout(context.Background(), timeout)
//	defer cancel()
//	_, err := tx.Exec(ctx, opts.SQL, opts.Args...)
//	return convertUserError(err)
//}
//
//func (d *Database) txQueryRow(tx pgx.Tx, opts Opts) errchk {
//	ctx, cancel := context.WithTimeout(context.Background(), timeout)
//	defer cancel()
//	err := tx.QueryRow(ctx, opts.SQL, opts.Args...).Scan(opts.Scan...)
//
//	return convertUserError(err)
//}
