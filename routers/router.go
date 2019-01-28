package routers

import (
	"github.com/gin-gonic/gin"
	"go-admin-starter/utils/config"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	conf := config.New()
	gin.SetMode(conf.Server.RunMode)

	initWebRouter(r)

	initApiRouter(r)

	return r
}
