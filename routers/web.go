package routers

import (
	"admin-server/controllers/web"
	"admin-server/pkg/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func initWebRouter(r *gin.Engine) *gin.Engine {
	r.LoadHTMLGlob("controllers/web/templates/*")

	r.GET("/welcome", web.Welcome)

	r.StaticFS("/upload", http.Dir(config.Conf.App.RuntimeRootPath + "upload"))
	r.Static("/docs", "docs/swagger")//docs

	return r
}
