// Package errorhandling implements some helper functions for handling errors.
// It has two main purposes: (a) to nicely errors for review when needed, and
// (b) to send to Sentry, if configured accordingly.
package errorhandling

import (
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4"
	"github.com/kaphos/webapp/internal/log"
)

// Check if an errchk object is throwing an errchk
// Wrapper function for noError, so that it makes more sense in the code
func Check(err error, errStr string) {
	HaveError(err, errStr)
}

// HaveError returns true if there are any errors (excluding db queries with no rows)
func HaveError(err error, errCode string) bool {
	if err == nil || err == ErrNoRows || err == ErrClientSide || err == pgx.ErrTxClosed {
		return false
	}

	log.Get("MAIN").Error(errCode + ": " + err.Error())

	if sentryInitialised {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			scope.SetTag("ref", errCode)
			sentry.CaptureException(err)
		})
	}

	return true

}
