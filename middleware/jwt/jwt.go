package jwt

import (
	"admin-server/pkg/app"
	"admin-server/pkg/e"
	"admin-server/pkg/redis"
	"admin-server/pkg/util"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			app.Response(c, e.INVALID_PARAMS, "缺少Token参数", nil)
			return
		}
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			app.Response(c, e.INVALID_PARAMS, "Token格式错误", nil)
			return
		}
		token = parts[1]

		isExist := redis.Master().Exists(token)
		if isExist.Val() != true {
			app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, "无效Token", nil)
			return
		}

		j := util.NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, "Token鉴权失败", nil)
			return
		} else if time.Now().Unix() > claims.ExpiresAt {
			app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, "Token已超时", nil)
			return
		}

		//继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)

		c.Next()
	}
}
