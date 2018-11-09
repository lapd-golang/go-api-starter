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

// @Summary 授权
// @Produce json
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {string} json "{"code":200,"data":{"token": "", "expires_at": 0},"message":"ok"}"
// @Router /auth [post]
func GetAuth(c *gin.Context) {
	valid := validation.Validation{}

	username := c.PostForm("username")
	password := c.PostForm("password")

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		app.MarkErrors(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	user := models.User{
		Username:username,
		Password:password,
	}
	user = user.CheckUser()

	if user.ID <= 0 {
		app.Response(c, e.ERROR_AUTH, "账号或密码错误", nil)
		return
	}

	j := util.NewJWT()
	tokenData, err := j.GenerateToken(user.ID, user.Username, c.Request.UserAgent())
	if err != nil {
		app.Response(c, e.ERROR_AUTH_TOKEN, "Token生成失败", nil)
		return
	}

	app.Response(c, e.SUCCESS, "ok", tokenData)
}