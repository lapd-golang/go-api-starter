package utils

import (
	"admin-server/utils/config"
	"admin-server/utils/redis"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//JWT 签名结构
type JWT struct {
	SigningKey []byte
}

//一些常量
var (
	TokenExpired     error = errors.New("Token is expired")
	TokenNotValidYet error = errors.New("Token not active yet")
	TokenMalformed   error = errors.New("That's not even a token")
	TokenInvalid     error = errors.New("Couldn't handle this token:")
)

//载荷，可以加一些自己需要的信息
type CustomClaims struct {
	ID        int    `json:"user_id"`
	Username  string `json:"username"`
	UserAgent string `json:"user_agent"`
	jwt.StandardClaims
}

//返回token结构
type TokenData struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

//新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(config.Conf.App.JwtSecret),
	}
}

//GenerateToken 生成一个token
func (j *JWT) GenerateToken(id int, username string, userAgent string) (TokenData, error) {
	nowTime := time.Now()
	lifeTime := 3 * time.Hour
	expireTime := nowTime.Add(lifeTime).Unix()

	//生成token
	claims := CustomClaims{
		id,
		username,
		userAgent,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "",
		},
	}

	tokenData := TokenData{}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(j.SigningKey)
	if err != nil {
		return tokenData, err
	}

	tokenData.TokenType = "Bearer"
	tokenData.AccessToken = accessToken
	tokenData.RefreshToken = EncodeMD5(accessToken)
	tokenData.ExpiresAt = expireTime

	//记录token到redis
	data, err := json.Marshal(tokenData)
	if err := redis.Master().Set(accessToken, data, lifeTime).Err(); err != nil {
		return tokenData, err
	}

	return tokenData, nil
}

//解析Tokne
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
