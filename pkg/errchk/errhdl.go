package errchk

import "github.com/kaphos/webapp/internal/errorhandling"

func Have(err error, errCode string) bool {
	return errorhandling.HaveError(err, errCode)
}
