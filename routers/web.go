package routers

import (
	"admin-server/controllers/web"
	"github.com/gin-gonic/gin"
)

func initWebRouter(r *gin.Engine) *gin.Engine {
	r.LoadHTMLGlob("controllers/web/templates/*")

	r.GET("/welcome", web.Welcome)

	r.Static("/static", "static")//静态资源目录，包含上传目录
	r.Static("/docs", "docs/swagger")//docs

	return r
}
