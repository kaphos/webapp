package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/kaphos/webapp/pkg/errchk"
)

func (d *Database) NewTransaction(ctx context.Context, spanName string, f func(tx pgx.Tx) error) error {
	ctx, span := d.tracer.Start(ctx, spanName)
	defer span.End()

	err := convertUserError(d.pool.BeginTxFunc(ctx, pgx.TxOptions{}, f))
	errchk.Check(err, "dbNewTx")

	return err
}
