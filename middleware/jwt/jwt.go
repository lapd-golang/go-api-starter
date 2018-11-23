package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-admin-starter/utils/app"
	"go-admin-starter/utils/e"
	"strings"
	"time"
)

//JWT签名结构
type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "kerlin"
)

//载荷
type Customclaims struct {
	ID        int    `json:"user_id"`
	Username  string `json:"username"`
	UserAgent string `json:"user_agent"`
	jwt.StandardClaims
}

//设置SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

//获取SignKey
func GetSignKey() string {
	return SignKey
}
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

//创建Token
func (j *JWT) CreateToken(claims Customclaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//解析token
func (j *JWT) ParseToken(tokenString string) (*Customclaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Customclaims{}, func(token *jwt.Token) (interface{}, error) {
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
	if claims, ok := token.Claims.(*Customclaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

//更新Token
func (j *JWT) RefreshToken(tokenString string, expiresAt int64) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &Customclaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*Customclaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = expiresAt
		return j.CreateToken(*claims)
	}

	return "", TokenInvalid
}

//JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			app.Response(c, e.INVALID_PARAMS, "缺少Token参数", nil)
			return
		}
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			app.Response(c, e.INVALID_PARAMS, "Token格式错误", nil)
			return
		}
		token = parts[1]

		//isExist := redis.Master().Exists(token)
		//if isExist.Val() != true {
		//	app.Response(c, e.ERROR_AUTH, "无效Token", nil)
		//	return
		//}

		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, "Token鉴权失败", nil)
			return
		} else if time.Now().Unix() > claims.ExpiresAt {
			app.Response(c, e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, "Token已超时", nil)
			return
		}

		//继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)

		c.Next()
	}
}
