package db

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/kaphos/webapp/internal/errorhandling"
	"os"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getDBConnStr(defaultUser, defaultPass string) string {
	dbUser := getenv("DB_USER", defaultUser)
	dbPass := getenv("DB_PASS", defaultPass)
	dbHost := getenv("DB_HOST", "127.0.0.1")
	dbPort := getenv("DB_PORT", "5432")
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s", dbUser, dbUser, dbPass, dbHost, dbPort)
}

// convertUserError returns a defined errchk (errorhandling.ErrClientSide) if the errchk code
// falls into a predefined set, that is due to user input errchk (e.g. duplicate).
func convertUserError(err error) error {
	if err != nil && err != pgx.ErrNoRows {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" || pgErr.Code == "23505" {
				return errorhandling.ErrClientSide
			}
		}
	} else if err == pgx.ErrNoRows {
		return errorhandling.ErrNoRows
	}
	return err
}
