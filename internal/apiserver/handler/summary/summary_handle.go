package summary

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/summary"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

// SummaryHandle 创建产品
func SummaryHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		logic := summary.NewSummaryLogic(c, svcCtx)
		resp, err := logic.Summary()
		response.HandleResponse(c, resp, err)
	}
}
