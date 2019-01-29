package utils

import (
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
	"go-admin-starter/utils/config"
)

var conf = config.New()

func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * conf.App.PageSize
	}

	return result
}