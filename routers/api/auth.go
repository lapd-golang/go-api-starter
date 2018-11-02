package api

import (
	"admin-server/models"
	"admin-server/pkg/app"
	"admin-server/pkg/e"
	"admin-server/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	valid := validation.Validation{}

	username := c.PostForm("username")
	password := c.PostForm("password")

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		app.MarkErrors(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, "请求参数错误", nil)
		return
	}

	isExist := models.CheckAuth(username, password)
	if !isExist {
		app.Response(c, e.ERROR_AUTH, "账号或密码错误", nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		app.Response(c, e.ERROR_AUTH_TOKEN, "Token生成失败", nil)
		return
	}

	app.Response(c, e.SUCCESS, "ok", map[string]string{
		"token": token,
	})
}