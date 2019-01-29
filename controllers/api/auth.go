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
	"golang.org/x/crypto/bcrypt"
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

//用户注册
func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	valid := validation.Validation{}
	valid.Required(username, "username").Message("用户名 不能为空")
	valid.Required(password, "password").Message("密码 不能为空")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, e.INVALID_PARAMS, valid.Errors[0].Message, nil)
		return
	}
	if password != confirmPassword {
		app.Response(c, e.INVALID_PARAMS, "两次输入密码不一致", nil)
		return
	}

	user := models.User{}
	if user.CheckExistByUsername(username) {
		app.Response(c, e.DATA_EXIST, "用户名 已注册", nil)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.Log.Errorf("Generate from password error:", err)
	}
	encodePW := string(hash) //保存在数据库的密码，虽然每次生成都不同，只需保存一份即可

	user.Username = username
	user.Password = encodePW
	id, err := user.Insert()
	if err != nil {
		app.Response(c, e.ERROR, err.Error(), nil)
		return
	}
	app.Response(c, e.SUCCESS, "ok", id)
	return
}

//授权登录
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

	expiresAt := time.Now().Add(lifeTime).Unix() //签名过期时间

	j := jwt.NewJWT()
	claims := jwt.Customclaims{
		user.ID,
		user.Username,
		c.Request.UserAgent(),
		user.Role,
		jwtgo.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    "kerlin", //签名发行者
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

//刷新accessToken
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

	//token是否存在于redis中
	isExist := redis.Master().Exists(accessToken)
	if isExist.Val() != true {
		app.Response(c, e.ERROR_AUTH, "无效Token", nil)
		return
	}

	//生成新token
	expiresAt := time.Now().Add(lifeTime).Unix() //签名过期时间
	newAccessToken, err := j.RefreshToken(accessToken, expiresAt)
	if err != nil {
		app.Response(c, e.ERROR_AUTH_TOKEN, "Token生成失败", nil)
		return
	}

	redis.Master().Del(accessToken) //移除旧token

	tokenData := TokenData{
		"Bearer",
		newAccessToken,
		utils.EncodeMD5(newAccessToken),
		expiresAt,
	}
	//记录token到redis
	data, err := json.Marshal(tokenData)
	if err := redis.Master().Set(newAccessToken, data, lifeTime).Err(); err != nil {
		utils.Log.Warn("recrod auth token to redis error: ", err)
	}

	app.Response(c, e.SUCCESS, "ok", tokenData)
}

/*------------------------------------后台管理员相关--------------------------------------------*/
//后台管理员授权登录
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

	expiresAt := time.Now().Add(lifeTime).Unix() //签名过期时间

	j := jwt.NewJWT()
	claims := jwt.Customclaims{
		user.ID,
		user.Username,
		c.Request.UserAgent(),
		user.Role,
		jwtgo.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    "kerlin", //签名发行者
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
