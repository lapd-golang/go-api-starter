package routers

import (
	"admin-server/utils/config"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(config.Conf.Server.RunMode)

	initWebRouter(r)

	initApiRouter(r)

	return r
}
