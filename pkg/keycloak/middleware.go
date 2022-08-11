package keycloak

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kaphos/webapp/internal/db"
	"github.com/kaphos/webapp/internal/log"
	"github.com/kaphos/webapp/pkg/repo"
)

//// Middleware for checking if a user has a valid Keycloak token.
//// Only checks if the JWT is signed and if it has a valid timestamp.
//// Does not perform any additional checks, such as whether the JWT
//// contains any specific roles. If that is needed,
//// BuildMiddlewareWithCheck should be used instead.
//var Middleware = repo.NewMiddleware(func(c *gin.Context) bool {
//	_, err := checkKeycloak(c)
//	return err == nil
//}, 401, "Unauthorised")

// BuildMiddlewareWithCheck returns a middleware that can be used to
// check if a user is authorised. It first checks for the JWT's
// validity. Then, it performs any further checks on the JWT
// claims as needed.
func (kc *Keycloak) BuildMiddlewareWithCheck(checkValid func(*gin.Context, *db.Database, jwt.MapClaims) bool) repo.Middleware {
	return repo.NewMiddleware(func(c *gin.Context) bool {
		claims, err := kc.checkKeycloak(c)
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
