package logs

import (
	logs "gitlab.intsig.net/textin-gateway/internal/apiserver/logic/logs/chart"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	logsType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"

	"github.com/gin-gonic/gin"
)

// AvgRespTimeLine 平均响应时间折线图
func AvgRespTimeLineHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {

	return func(c *gin.Context) {

		var req logsType.GraphReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := logs.NewAvgRespTimeLineLogic(c, svcCtx)
		resp, err := logic.AvgRespTimeLine(&req)
		response.HandleResponse(c, resp, err)
	}
}
