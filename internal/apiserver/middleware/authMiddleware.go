package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	authType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/pkg/jwtgen"
	"gorm.io/gorm"
)

type Dispatcher struct {
	pgClient *gorm.DB // postgres 连接池
}

func AuthMiddleware(c *gin.Context) {
	// TODO generate middleware implement function, delete after code implementation
	claims, _ := verifyJWT(c)
	c.Set("user", claims)
	c.Next()
}

func verifyJWT(c *gin.Context) (userInfo authType.UserInfo, err error) {
	var (
		TokenExpired error = errors.New("Token is expired")
		// TokenNotValidYet error = errors.New("Token not active yet")
		// TokenMalformed   error = errors.New("That's not even a token")
		// TokenInvalid     error = errors.New("Couldn't handle this token:")
	)

	authorizationHeader := c.Request.Header.Get("Authorization")
	if authorizationHeader == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"code": 401,
			"msg":  "Authorization header is required",
		})
		return
	}
	token := strings.TrimPrefix(authorizationHeader, "Bearer ")
	if token == authorizationHeader {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "Authorization header must be in the format of Bearer {token}",
		})
		return
	}

	log.Print("get token: ", token)

	// parseToken 解析token包含的信息
	claims, err := jwtgen.ParseToken(token)
	if err != nil {
		if err == TokenExpired {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "授权已过期",
			})
			c.Abort()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  err.Error(),
		})
		c.Abort()
		return
	}
	userInfo = claims.UserInfo
	if err != userInfo.GetRole() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户获取失败",
		})
		c.Abort()
		return
	}
	return
}
