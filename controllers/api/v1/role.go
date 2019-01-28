package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"go-admin-starter/models"
	"go-admin-starter/utils/app"
	"go-admin-starter/utils/e"
)

func AddCasbin(c *gin.Context) {
	rolename := c.PostForm("rolename")
	path := c.PostForm("path")
	method := c.PostForm("method")

	valid := validation.Validation{}
	valid.Required(rolename, "rolename").Message("角色名称不能为空")
	valid.Required(path, "path").Message("路径不能为空")
	valid.Required(method, "method").Message("请求方法不能为空")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	ptype := "p"
	casbin := models.CasbinModel{
		Ptype:    ptype,
		RoleName: rolename,
		Path:     path,
		Method:   method,
	}
	isok := casbin.AddCasbin(casbin)
	if isok {
		app.Response(c, e.SUCCESS, "保存成功", nil)
		return
	} else {
		app.Response(c, e.ERROR, "保存失败", nil)
		return
	}
}
