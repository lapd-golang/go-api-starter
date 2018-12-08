package api

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-admin-starter/middleware/jwt"
	"go-admin-starter/models"
	"go-admin-starter/utils"
	"go-admin-starter/utils/app"
	"go-admin-starter/utils/e"
	"go-admin-starter/utils/redis"
	"time"
)

var lifeTime = 3 * time.Hour

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

//返回token结构
type TokenData struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
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

	expiresAt := time.Now().Add(lifeTime).Unix()//签名过期时间

	j := jwt.NewJWT()
	claims := jwt.Customclaims{
		user.ID,
		user.Username,
		c.Request.UserAgent(),
		user.Role,
		jwtgo.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    "kerlin",//签名发行者
		},
	}
	accessToken, err := j.CreateToken(claims)
	if err != nil {
		app.Response(c, e.ERROR_AUTH_TOKEN, "Token生成失败", nil)
		return
	}

	tokenData := TokenData{
		"Bearer",
		accessToken,
		utils.EncodeMD5(accessToken),
		expiresAt,
	}
	//记录token到redis
	data, err := json.Marshal(tokenData)
	if err := redis.Master().Set(accessToken, data, lifeTime).Err(); err != nil {
		utils.Log.Warn("recrod auth token to redis error: ", err)
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

	//解析accessToken，验证过期
	j := jwt.NewJWT()
	claims, err := j.ParseToken(accessToken)
	if err != nil {
		app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, "Token鉴权失败", err)
		return
	} else if time.Now().Unix() > claims.ExpiresAt {
		app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, "Token已超时", nil)
		return
	}

	//判断是否相同UserAgent；验证refreshToken
	if claims.UserAgent != c.Request.UserAgent() || refreshToken != utils.EncodeMD5(accessToken) {
		app.Response(c, e.ERROR_AUTH, "无效Token", nil)
		return
	}

	//生成新token
	expiresAt := time.Now().Add(lifeTime).Unix()//签名过期时间
	newAccessToken, err := j.RefreshToken(accessToken, expiresAt)
	if err != nil {
		app.Response(c, e.ERROR_AUTH_TOKEN, "Token生成失败", nil)
		return
	}

	redis.Master().Del(accessToken)//移除旧token

	tokenData := TokenData{
		"Bearer",
		newAccessToken,
		utils.EncodeMD5(newAccessToken),
		expiresAt,
	}
	//记录token到redis
	data, err := json.Marshal(tokenData)
	if err := redis.Master().Set(accessToken, data, lifeTime).Err(); err != nil {
		utils.Log.Warn("recrod auth token to redis error: ", err)
	}

	app.Response(c, e.SUCCESS, "ok", tokenData)
}

/*------------------------------------后台管理员相关--------------------------------------------*/
// @Summary 后台管理员授权
// @Produce json
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {string} json "{"code":200,"data":{"token_type": "", "access_token": "", "refresh_token": "", "expires_at": 0},"message":"ok"}"
// @Router /admin/auth [post]
func AdminGetAuth(c *gin.Context) {
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

	user := models.AdminUser{
		Username: username,
		Password: password,
	}
	user = user.CheckUser()

	if user.ID <= 0 {
		app.Response(c, e.ERROR_AUTH, "账号或密码错误", nil)
		return
	}

	expiresAt := time.Now().Add(lifeTime).Unix()//签名过期时间

	j := jwt.NewJWT()
	claims := jwt.Customclaims{
		user.ID,
		user.Username,
		c.Request.UserAgent(),
		user.Role,
		jwtgo.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    "kerlin",//签名发行者
		},
	}
	accessToken, err := j.CreateToken(claims)
	if err != nil {
		app.Response(c, e.ERROR_AUTH_TOKEN, "Token生成失败", nil)
		return
	}

	tokenData := TokenData{
		"Bearer",
		accessToken,
		utils.EncodeMD5(accessToken),
		expiresAt,
	}
	//记录token到redis
	data, err := json.Marshal(tokenData)
	if err := redis.Master().Set(accessToken, data, lifeTime).Err(); err != nil {
		utils.Log.Warn("recrod auth token to redis error: ", err)
	}

	app.Response(c, e.SUCCESS, "ok", tokenData)
}
