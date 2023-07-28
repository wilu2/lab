package auth

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/handler/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoute(e *gin.Engine, svcCtx *svc.ServiceContext) {
	g := e.Group("/gateway/v1")
	g.POST("/login_verify", auth.LoginVerifyHandle(svcCtx))
	g.POST("/token_verify", auth.TokenVerifyHandle(svcCtx))
}
