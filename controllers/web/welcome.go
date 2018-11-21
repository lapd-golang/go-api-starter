package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Welcome(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome.tmpl", gin.H{
		"title": "Welcome",
		"content": "Welcome!",
	})
}
