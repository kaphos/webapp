package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/kaphos/webapp/pkg/errchk"
)

func (d *Database) NewTransaction(ctx context.Context, spanName string, f func(tx pgx.Tx) error) error {
	ctx, span := d.tracer.Start(ctx, spanName)
	defer span.End()

	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{})
	if errchk.HaveError(err, "dbNewTx0") {
		return err
	}
	err = convertUserError(f(tx))
	errchk.Check(err, "dbNewTx1")

	return err
}
