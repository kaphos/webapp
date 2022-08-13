package keycloak

import (
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/pkg/db"
	"github.com/kaphos/webapp/pkg/errchk"
	"github.com/kaphos/webapp/pkg/middleware"
)

type Keycloak struct {
	key        *rsa.PublicKey
	DB         *db.Database
	Middleware middleware.Middleware
}

func New(publicKey string, database *db.Database) (Keycloak, error) {
	kc := Keycloak{DB: database}

	if err := kc.parsePublicKey(publicKey); errchk.HaveError(err, "newKC") {
		return kc, err
	}

	kc.Middleware = middleware.NewAuth(func(c *gin.Context) bool {
		_, err := kc.Verify(c)
		return err == nil
	})

	return kc, nil
}
