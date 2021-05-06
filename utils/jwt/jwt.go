package jwt

import (
	"fmt"
	"goblog/utils/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	// ID         uint32
	User       string
	JWTVersion int
	jwt.StandardClaims
}

func NewJWT(c *CustomClaims) (string, error) {
	// 设置版本号，有更改时可以使之前签发的jwt失效
	c.JWTVersion = config.JwtVersion()
	// 设置过期时间
	c.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
		Issuer:    config.JwtIssuer(),
	}
	// 设置加密算法
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 生成jwt字符串
	token, err := tokenClaims.SignedString([]byte(config.JwtKey()))
	if err != nil {
		return "", fmt.Errorf("生成JWT失败")
	}
	return token, nil
}

func ParseJWT(token string) (*CustomClaims, error) {
	// 解析jwt到自定义结构
	tokenClaims, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtKey()), nil
	})
	if err != nil || tokenClaims == nil {
		return nil, fmt.Errorf("解析JWT失败")
	}
	// 校验token是否有效
	if claims, ok := tokenClaims.Claims.(*CustomClaims); ok && tokenClaims.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("JWT无效或已失效")
	}
}
