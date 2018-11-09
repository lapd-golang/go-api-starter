package api

import (
	"admin-server/models"
	"admin-server/pkg/app"
	"admin-server/pkg/e"
	"admin-server/pkg/redis"
	"admin-server/pkg/util"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"time"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary 授权
// @Produce json
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {string} json "{"code":200,"data":{"token_type": "", "access_token": "", "refresh_token": "", "expires_at": 0},"message":"ok"}"
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
		Username: username,
		Password: password,
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


// @Summary 刷新accessToken
// @Produce json
// @Param access_token formData string true "access_token"
// @Param refresh_token formData string true "refresh_token"
// @Success 200 {string} json "{"code":200,"data":{"token_type": "", "access_token": "", "refresh_token": "", "expires_at": 0},"message":"ok"}"
// @Router /refreshToken [post]
func RefreshToken(c *gin.Context) {
	accessToken := c.PostForm("access_token")
	refreshToken := c.PostForm("refresh_token")

	valid := validation.Validation{}
	valid.Required(accessToken, "access_token").Message("access_token 不能为空")
	valid.Required(refreshToken, "refresh_token").Message("refresh_token 不能为空")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}

	//验证accessToken是否存在
	dbToken := redis.Master().Get(accessToken)
	if dbToken.Val() == "" {
		app.Response(c, e.ERROR_AUTH, "无效Token", nil)
		return
	}

	//解析accessToken，验证过期
	j := util.NewJWT()
	claims, err := j.ParseToken(accessToken)
	if err != nil {
		app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, "Token鉴权失败", nil)
		return
	} else if time.Now().Unix() > claims.ExpiresAt {
		app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, "Token已超时", nil)
		return
	}

	//判断是否相同UserAgent；验证refreshToken
	var t util.TokenData
	json.Unmarshal([]byte(dbToken.Val()), &t)
	if claims.UserAgent != c.Request.UserAgent() || refreshToken != t.RefreshToken {
		app.Response(c, e.ERROR_AUTH, "无效Token", nil)
		return
	}

	//生成新token
	tokenData, err := j.GenerateToken(claims.ID, claims.Username, c.Request.UserAgent())
	if err != nil {
		app.Response(c, e.ERROR_AUTH_TOKEN, "Token生成失败", nil)
		return
	}

	app.Response(c, e.SUCCESS, "ok", tokenData)
}
