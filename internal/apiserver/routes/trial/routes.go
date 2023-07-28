package trial

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/handler/trial"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/middleware"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

func RegisterTrialRoute(e *gin.Engine, svcCtx *svc.ServiceContext) {
	g := e.Group("/gateway/v1")
	g.Use(middleware.AuthMiddleware)
	g.GET("/trial_services", trial.ListTrialServicesHandle(svcCtx))
	g.POST("/trials", trial.CreateTrialHandle(svcCtx))
	g.POST("/trial/schema_validate", trial.SchemaValidateHandle(svcCtx))
}
