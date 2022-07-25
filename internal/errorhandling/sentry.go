package errorhandling

import (
	"github.com/getsentry/sentry-go"
	"github.com/kaphos/webapp/internal/log"
	"os"
)

var sentryLogger = log.Get("SENTRY")
var sentryInitialised = false

// InitSentry initialises Sentry, given `SENTRY_URL` is provided as an environment variable.
// Fails gracefully if it is not. Used to forward errors to a centralised platform, and so
// the `SENTRY_URL` env var should only be included in staging/production.
func InitSentry() {
	sentryUrl := os.Getenv("SENTRY_URL")
	if sentryUrl == "" {
		return
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_URL"),
	}); err != nil {
		sentryLogger.Error("Error initialising Sentry: " + err.Error())
	}

	sentryLogger.Info("Initialised Sentry.")
	sentryInitialised = true
}
