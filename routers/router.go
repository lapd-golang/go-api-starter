package routers

import (
	"github.com/gin-gonic/gin"
	"go-admin-starter/utils/config"
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
