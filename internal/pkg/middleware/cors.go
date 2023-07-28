package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	maxAge = 12
)

// Cors add cors headers.
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"http://global-test.textin.com", "https://global-test.textin.com",
			"https://textin.com", "*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "noMes"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: maxAge * time.Hour,
	})
}
