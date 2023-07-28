package logs

import (
	logs "gitlab.intsig.net/textin-gateway/internal/apiserver/logic/logs/chart"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	logsType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"

	"github.com/gin-gonic/gin"
)

// SummaryHandle 创建产品
func ConcurGraphLogsHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {

	return func(c *gin.Context) {

		var req logsType.GraphReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := logs.NewConcurLineLogic(c, svcCtx)
		resp, err := logic.ConcurLine(&req)
		response.HandleResponse(c, resp, err)
	}
}
