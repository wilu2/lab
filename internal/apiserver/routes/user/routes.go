package user

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/handler/user"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/middleware"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoute(e *gin.Engine, svcCtx *svc.ServiceContext) {
	g := e.Group("/gateway/v1")
	g.Use(middleware.AuthMiddleware)
	g.POST("/users", user.CreateUserHandle(svcCtx))
	g.GET("/users/:id", user.GetUserHandle(svcCtx))
	g.POST("/users/:id/update", user.UpdateUserHandle(svcCtx))
	g.POST("/users/:id/delete", user.DeleteUserHandle(svcCtx))
	g.GET("/users", user.ListUserHandle(svcCtx))
	g.POST("/reset_password", user.ResetPasswordHandle(svcCtx))
}
