package logs

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/logs"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	logsType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"

	"github.com/gin-gonic/gin"
)

// SummaryHandle 创建产品
func ListServiceHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {

	return func(c *gin.Context) {

		var req logsType.ServiceListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := logs.NewListServiceLogic(c, svcCtx)
		resp, err := logic.ListService(&req)
		response.HandleResponse(c, resp, err)
	}
}
