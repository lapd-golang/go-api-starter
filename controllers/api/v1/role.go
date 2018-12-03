package v1

import (
	"github.com/gin-gonic/gin"
	"go-admin-starter/models"
	"net/http"
)

func AddCasbin(c *gin.Context) {
	rolename := c.PostForm("rolename")
	path := c.PostForm("path")
	method := c.PostForm("method")
	ptype := "p"
	casbin := models.CasbinModel{
		Ptype:    ptype,
		RoleName: rolename,
		Path:     path,
		Method:   method,
	}
	isok := casbin.AddCasbin(casbin)
	if isok {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "保存成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "保存失败",
		})
	}
}
