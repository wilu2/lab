package logs

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/handler/logs"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/middleware"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

func RegisterLogsRoute(e *gin.Engine, svcCtx *svc.ServiceContext) {
	g := e.Group("/gateway/v1/logs")
	g.Use(middleware.AuthMiddleware)
	g.GET("/list_service", logs.ListServiceHandle(svcCtx))
	g.GET("/list_channel", logs.ListChannelHandle(svcCtx))
	g.GET("/list_application", logs.ListApplicationHandle(svcCtx))
	g.GET("/list_version", logs.ListVersionHandle(svcCtx))

	g.POST("/query", logs.QueryLogsHandle(svcCtx))

	g.POST("/sum_request_line", logs.SumRequestLineHandle(svcCtx))
	g.POST("/concur_line", logs.ConcurGraphLogsHandle(svcCtx))
	g.POST("/dist_column", logs.DistColumnHandle(svcCtx))
	g.POST("/avg_resp_time_line", logs.AvgRespTimeLineHandle(svcCtx))

	g.GET("/get_request", logs.GetRequestBodyHandle(svcCtx))
	g.GET("/get_response", logs.GetResponseBodyHandle(svcCtx))
}
