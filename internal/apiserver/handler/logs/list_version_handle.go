package logs

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/logs"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	logsType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"

	"github.com/gin-gonic/gin"
)

// SummaryHandle 创建产品
func ListVersionHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {

	return func(c *gin.Context) {

		var req logsType.VersionListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := logs.NewListVersionLogic(c, svcCtx)
		resp, err := logic.ListVersion(&req)
		response.HandleResponse(c, resp, err)
	}
}
