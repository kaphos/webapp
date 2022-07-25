package errorhandling

import "fmt"

var (
	ErrClientSide = fmt.Errorf("errchk due to bad request")
	ErrNoRows     = fmt.Errorf("no rows in result set")
)
