// Package db provides a Database struct that builds on the pgx package,
// while providing standardised logging, errchk handling, and instrumentation.
package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaphos/webapp/internal/telemetry"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"math/rand"
	"time"

	"github.com/kaphos/webapp/internal/log"
)

// Database - used to connect to the database
type Database struct {
	pool   *pgxpool.Pool
	tracer trace.Tracer
	logger *zap.Logger
}

// NewDB initialises a new Database object, creating a Database pool and setting up logging
// and telemetry.
func NewDB(appName, defaultUser, defaultPass string, maxConns int32) (*Database, error) {
	rand.Seed(time.Now().UTC().UnixNano()) // set rand seed just in case. useful for testing.

	d := Database{
		logger: log.Get("DB"),
		tracer: telemetry.NewTracer(appName, "database"),
	}

	config, err := pgxpool.ParseConfig(getDBConnStr(defaultUser, defaultPass))
	if err != nil {
		d.logger.Error("Unable to parse database config: " + err.Error())
		return &Database{}, err
	}

	config.MaxConns = maxConns

	d.pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		d.logger.Error("Unable to connect to db: " + err.Error())
		return &Database{}, err
	}

	d.logger.Info("Connected to database.")

	return &d, nil
}

func (d *Database) Healthcheck(ctx context.Context) error {
	return d.pool.Ping(ctx)
}
