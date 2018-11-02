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
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			app.Response(c, e.INVALID_PARAMS, "请求参数错误", nil)
			return
		}
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			app.Response(c, e.INVALID_PARAMS, "请求参数错误", nil)
			return
		}
		token = parts[1]

		claims, err := util.ParseToken(token)
		if err != nil {
			app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, "Token鉴权失败", nil)
			return
		} else if time.Now().Unix() > claims.ExpiresAt {
			app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, "Token已超时", nil)
			return
		}

		c.Next()
	}
}
