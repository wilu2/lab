package summary

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/summary"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	summaryType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/summary"

	"github.com/gin-gonic/gin"
)

// SummaryHandle 创建产品
func SummaryGraphHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {

	return func(c *gin.Context) {

		var req summaryType.SummaryReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := summary.NewSummaryGraphLogic(c, svcCtx)
		resp, err := logic.SummaryGraph(&req)
		response.HandleResponse(c, resp, err)
	}
}
