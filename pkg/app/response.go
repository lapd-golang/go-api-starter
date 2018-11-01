package app

import (
	"admin-server/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"message": e.GetMsg(code),
		"data": data,
	})

	c.Abort()
}
