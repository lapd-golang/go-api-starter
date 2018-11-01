package jwt

import (
	"admin-server/pkg/app"
	"admin-server/pkg/e"
	"admin-server/pkg/util"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data interface{}

		token := c.Request.Header.Get("Authorization")
		if token == "" {
			app.Response(c, e.INVALID_PARAMS, data)
			return
		}
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			app.Response(c, e.INVALID_PARAMS, data)
			return
		}
		token = parts[1]

		claims, err := util.ParseToken(token)
		if err != nil {
			app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, data)
			return
		} else if time.Now().Unix() > claims.ExpiresAt {
			app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, data)
			return
		}

		c.Next()
	}
}
