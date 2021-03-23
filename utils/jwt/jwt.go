package jwt

import (
	"fmt"
	"goblog/utils/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	ID         uint32
	GithubID   uint32
	GithubUser string
	GithubName string
	JWTVersion int
	jwt.StandardClaims
}

func NewJWT(c *CustomClaims) (token string, err error) {
	c.JWTVersion = 1
	// set expired time of jwt
	c.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
		Issuer:    config.JwtIssuer(),
	}
	// set signing method of jwt
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// generate jwt string
	token, err = tokenClaims.SignedString([]byte(config.JwtKey()))
	return
}

func ParseJWT(token string) (*CustomClaims, error) {
	// parse the token string to custom claim struct
	tokenClaims, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtKey()), nil
	})
	if err != nil || tokenClaims == nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*CustomClaims); ok && tokenClaims.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("failed to valid the jwt")
	}
}
