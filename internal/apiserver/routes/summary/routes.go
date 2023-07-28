package summary

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/handler/summary"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/middleware"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

func RegisterSummaryRoute(e *gin.Engine, svcCtx *svc.ServiceContext) {
	g := e.Group("/gateway/v1")
	g.Use(middleware.AuthMiddleware)
	g.GET("/summary", summary.SummaryHandle(svcCtx))
	g.POST("/summary_graph", summary.SummaryGraphHandle(svcCtx))
}
