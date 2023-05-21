package db

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/kaphos/webapp/pkg/errchk"
	"os"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getDBConnStr(defaultUser, defaultPass string) string {
	dbUser := getEnv("DB_USER", defaultUser)
	dbName := getEnv("DB_NAME", defaultUser)
	dbPass := getEnv("DB_PASS", defaultPass)
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "5432")
	unixSocket := getEnv("INSTANCE_UNIX_SOCKET", "")
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
