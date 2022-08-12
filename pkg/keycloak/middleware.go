package keycloak

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kaphos/webapp/internal/log"
	"github.com/kaphos/webapp/pkg/db"
	"github.com/kaphos/webapp/pkg/errchk"
	"github.com/kaphos/webapp/pkg/repo"
)

// MiddlewareWithCheckFn returns a middleware that can be used to
// check if a user is authorised. It first checks for the JWT's
// validity. Then, it performs any further checks on the JWT
// claims as needed.
func (kc *Keycloak) MiddlewareWithCheckFn(checkValid func(*gin.Context, *db.Database, jwt.MapClaims) bool) repo.Middleware {
	return repo.NewMiddleware(func(c *gin.Context) bool {
		claims, err := kc.Verify(c)
		if err != nil {
			return false
		}

		if !checkValid(c, kc.DB, claims) {
			log.Get("KC").Debug("Did not pass checkValid function")
			return false
		}

		return true
	}, 401, "Unauthorised")
}

// MiddlewareWithIDCheck returns a middleware that can be used to
// check if a user is authorised. Takes in an SQL query that should
// expect 1 parameter (where a keycloak sub is passed in) and returns
// the ID of the user. For example, "SELECT id FROM users WHERE kc_sub = $1".
func (kc *Keycloak) MiddlewareWithIDCheck(sqlQuery string) repo.Middleware {
	return repo.NewMiddleware(func(c *gin.Context) bool {
		claims, err := kc.Verify(c)
		if err != nil {
			return false
		}

		var id int

		sub := claims["sub"].(string)
		err = kc.DB.QueryRow("kc-verify", c.Request.Context(), sqlQuery, sub).Scan(&id)
		if errchk.HaveError(err, "kcVerify0") {
			return false
		}

		c.Set("user-id", id) // note that id can be 0, if there's no match

		return true
	}, 401, "Unauthorised")
}
