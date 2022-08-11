package keycloak

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kaphos/webapp/internal/log"
	"github.com/kaphos/webapp/pkg/errchk"
	"net/http"
	"strings"
)

//var key *rsa.PublicKey

func (kc *Keycloak) parsePublicKey(publicKey string) error {
	var err error
	kc.key, err = jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	errchk.Check(err, "kcParsePubKey")
	return err
}

func (kc *Keycloak) extractKeyFromToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("error with jwt")
	}
	return kc.key, nil
}

// VerifyToken verifies that a Bearer token comes from an authorised
// Keycloak instance, and that the JWT it contains is valid.
func (kc *Keycloak) VerifyToken(authToken string) (jwt.MapClaims, error) {
	if kc.key == nil {
		err := fmt.Errorf("public key not initialised")
		errchk.Check(err, "kcVerifyToken")
		return nil, err
	}

	tokenSplit := strings.Split(authToken, " ")
	if len(tokenSplit) != 2 {
		return nil, fmt.Errorf("token string incorrect length")
	}

	token := tokenSplit[1]
	parsedToken, err := jwt.Parse(token, kc.extractKeyFromToken, jwt.WithoutClaimsValidation())
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return parsedToken.Claims.(jwt.MapClaims), nil
}

func (kc *Keycloak) checkKeycloak(c *gin.Context) (jwt.MapClaims, error) {
	authToken := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authToken, "Bearer ") {
		c.AbortWithStatus(http.StatusUnauthorized)
		log.Get("KC").Debug("Token missing Bearer")
		return nil, fmt.Errorf("invalid token")
	}

	claims, err := kc.VerifyToken(authToken)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		log.Get("KC").Debug("Token unverified: " + err.Error())
		return nil, fmt.Errorf("invalid token")
	}

	c.Set("kc-id", claims["sub"])
	c.Set("kc-roles", claims["realm_access"].(map[string]interface{})["roles"])

	return claims, nil
}
