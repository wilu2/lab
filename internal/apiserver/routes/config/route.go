package config

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/handler/config"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/middleware"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

func RegisterConfigRoute(e *gin.Engine, svcCtx *svc.ServiceContext) {
	g := e.Group("/gateway/v1")
	g.Use(middleware.AuthMiddleware)
	g.POST("/config", config.UpdateConfigHandle(svcCtx))
}

func RegisterConfigRouteWithoutAuth(e *gin.Engine, svcCtx *svc.ServiceContext) {
	g := e.Group("/gateway/v1")
	g.GET("/config", config.GetConfigHandle(svcCtx))
}
