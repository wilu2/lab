package logs

import (
	logs "gitlab.intsig.net/textin-gateway/internal/apiserver/logic/logs/chart"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	logsType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"

	"github.com/gin-gonic/gin"
)

// DistColumn 状态码分布柱状图
func DistColumnHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {

	return func(c *gin.Context) {

		var req logsType.GraphReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := logs.NewDistColumnLogic(c, svcCtx)
		resp, err := logic.DistColumn(&req)
		response.HandleResponse(c, resp, err)
	}
}
