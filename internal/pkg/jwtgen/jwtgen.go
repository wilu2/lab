package jwtgen

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	authType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
)

// 自定义Claims
type CustomClaims struct {
	authType.UserInfo
	jwt.StandardClaims
}

// 生成token
func GetToken(userInfo authType.UserInfo) (tokenString string) {
	var (
		timeOut = viper.GetDuration("jwt.timeout")
		key     = viper.GetString("jwt.key")
	)
	fmt.Println(timeOut)
	customClaims := &CustomClaims{
		UserInfo: userInfo,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(timeOut).Unix(), // 过期时间，必须设置
			Issuer:    "jerry",                        // 非必须，也可以填充用户名，
		},
	}
	//采用HMAC SHA256加密算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("token: %v\n", tokenString)
	return
}

// 解析token
func ParseToken(tokenString string) (*CustomClaims, error) {
	var (
		key = viper.GetString("jwt.key")
	)
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	// fmt.Println(token.Valid)
	// if token.Valid {
	// 	if claims, ok := token.Claims.(*CustomClaims); ok {
	// 		return claims, nil
	// 	}
	// }
	// return nil, err
	if token != nil {
		claims, ok := token.Claims.(*CustomClaims)
		if ok && token.Valid {
			return claims, nil
		}
	}

	return nil, err
}
