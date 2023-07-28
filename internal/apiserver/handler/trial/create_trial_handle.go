package trial

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/trial"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	trialType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/trial"

	"github.com/gin-gonic/gin"
)

// CreateTrialHandle 删除产品
func CreateTrialHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req trialType.TrialRouteDef
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := trial.NewCreateTrialLogic(c, svcCtx)
		resp, err := logic.CreateTrial(&req)
		response.HandleResponse(c, resp, err)
	}
}
