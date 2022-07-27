package keycloak

import (
	"crypto/rsa"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kaphos/webapp/internal/log"
	"github.com/kaphos/webapp/pkg/errchk"
	"github.com/kaphos/webapp/pkg/repo"
	"net/http"
	"strings"
)

var key *rsa.PublicKey

// ParsePublicKey parses a PEM-formatted public key. Should be called
// once, at server initialisation and before any requests are handled.
// Used to verify that a given JWT comes from our Keycloak client.
func ParsePublicKey(publicKey string) error {
	var err error
	key, err = jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	errchk.Check(err, "kcParsePubKey")
	return err
}

func extractKeyFromToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("error with jwt")
	}
	return key, nil
}

// VerifyToken verifies that a Bearer token comes from an authorised
// Keycloak instance, and that the JWT it contains is valid.
func VerifyToken(authToken string) (jwt.MapClaims, error) {
	if key == nil {
		err := fmt.Errorf("public key not initialised")
		errchk.Check(err, "kcVerifyToken")
		return nil, err
	}

	tokenSplit := strings.Split(authToken, " ")
	if len(tokenSplit) != 2 {
		return nil, fmt.Errorf("token string incorrect length")
	}

	token := tokenSplit[1]
	parsedToken, err := jwt.Parse(token, extractKeyFromToken, jwt.WithoutClaimsValidation())
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return parsedToken.Claims.(jwt.MapClaims), nil
}

func checkKeycloak(c *gin.Context) (jwt.MapClaims, error) {
	authToken := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authToken, "Bearer ") {
		c.AbortWithStatus(http.StatusUnauthorized)
		log.Get("KC").Debug("Token missing Bearer")
		return nil, fmt.Errorf("invalid token")
	}

	claims, err := VerifyToken(authToken)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		log.Get("KC").Debug("Token unverified: " + err.Error())
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// Middleware for checking if a user has a valid Keycloak token.
// Only checks if the JWT is signed and if it has a valid timestamp.
// Does not perform any additional checks, such as whether the JWT
// contains any specific roles. If that is needed,
// BuildMiddlewareWithCheck should be used instead.
var Middleware = repo.NewMiddleware(func(c *gin.Context) bool {
	_, err := checkKeycloak(c)
	return err == nil
}, 401, "Unauthorised")

// BuildMiddlewareWithCheck returns a middleware that can be used to
// check if a user is authorised. It first checks for the JWT's
// validity. Then, it performs any further checks on the JWT
// claims as needed.
func BuildMiddlewareWithCheck(checkValid func(jwt.MapClaims) bool) repo.Middleware {
	return repo.NewMiddleware(func(c *gin.Context) bool {
		claims, err := checkKeycloak(c)
		if err != nil {
			return false
		}

		if !checkValid(claims) {
			log.Get("KC").Debug("Did not pass checkValid function")
			return false
		}

		return true
	}, 401, "Unauthorised")
}
