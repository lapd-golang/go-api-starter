package authrole

import (
	"github.com/gin-gonic/gin"
	"go-admin-starter/middleware/jwt"
	"go-admin-starter/models"
	"go-admin-starter/utils/app"
	"go-admin-starter/utils/e"
)

//权限检查中间件
func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		//根据上下文获取载荷claims 从claims获得role
		claims := c.MustGet("claims").(*jwt.Customclaims)
		role := claims.Role
		enforcer := models.Casbin()
		//检查权限
		res, err := enforcer.EnforceSafe(role, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			app.Response(c, e.ERROR_AUTH_ROLE, "检查权限错误: "+err.Error(), nil)
			return
		}
		if res {
			c.Next()
		} else {
			app.Response(c, e.ERROR_AUTH_ROLE, "很抱歉您没有此权限", nil)
			return
		}
	}
}
