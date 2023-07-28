package trial

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/trial"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	trialType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/trial"

	"github.com/gin-gonic/gin"
)

// ListTrialServicesHandle 创建产品
func ListTrialServicesHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req trialType.ListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := trial.NewListTrialServiceLogic(c, svcCtx)
		resp, err := logic.ListTrialService(&req)
		response.HandleResponse(c, resp, err)
	}
}
