package logs

import (
	logs "gitlab.intsig.net/textin-gateway/internal/apiserver/logic/logs/chart"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	logsType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"

	"github.com/gin-gonic/gin"
)

// SumRequestLine 总请求量折线图
func SumRequestLineHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {

	return func(c *gin.Context) {

		var req logsType.GraphReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := logs.NewSumRequestLineLogic(c, svcCtx)
		resp, err := logic.SumRequestLine(&req)
		response.HandleResponse(c, resp, err)
	}
}
