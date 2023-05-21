package db

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/kaphos/webapp/pkg/errchk"
	"github.com/kaphos/webapp/pkg/utils"
)

func getDBConnStr(defaultUser, defaultPass string) string {
	dbUser := utils.GetEnv("DB_USER", defaultUser)
	dbName := utils.GetEnv("DB_NAME", defaultUser)
	dbPass := utils.GetEnv("DB_PASS", defaultPass)
	dbHost := utils.GetEnv("DB_HOST", "127.0.0.1")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	unixSocket := utils.GetEnv("INSTANCE_UNIX_SOCKET", "")
	if unixSocket != "" {
		return fmt.Sprintf("dbname=%s user=%s password=%s host=%s", dbName, dbUser, dbPass, unixSocket)
	}
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s", dbUser, dbUser, dbPass, dbHost, dbPort)
}

// convertUserError returns a defined errchk (errchk.ErrClientSide) if the errchk code
// falls into a predefined set, that is due to user input errchk (e.g. duplicate).
func convertUserError(err error) error {
	if err == nil || err == pgx.ErrNoRows || err == pgx.ErrTxClosed {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23505") {
		return errchk.ErrClientSide
	}

	return err
}
