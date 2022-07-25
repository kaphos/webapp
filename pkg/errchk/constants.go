package errchk

import "fmt"

var (
	ErrClientSide = fmt.Errorf("bad request")
	ErrNoRows     = fmt.Errorf("no rows in result set")
)
