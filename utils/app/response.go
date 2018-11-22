package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"message": message,
		"data": data,
	})

	c.Abort()
}
