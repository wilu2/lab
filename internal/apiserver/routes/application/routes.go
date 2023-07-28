// Code generated by goctl. DO NOT EDIT.
package application

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/handler/application"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/middleware"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

func RegisterApplicationRoute(e *gin.Engine, svcCtx *svc.ServiceContext) {
	g := e.Group("/gateway/v1")
	g.Use(middleware.AuthMiddleware)
	g.GET("/applications", application.ListApplicationHandle(svcCtx))
	g.POST("/applications", application.CreateApplicationHandle(svcCtx))
	g.GET("/applications/:id", application.GetApplicationHandle(svcCtx))
	g.POST("/applications/:id/update", application.UpdateApplicationHandle(svcCtx))
	g.POST("/applications/:id/update_status", application.UpdateStatusApplicationHandle(svcCtx))
	g.POST("/applications/:id/delete", application.DeleteApplicationHandle(svcCtx))

}
