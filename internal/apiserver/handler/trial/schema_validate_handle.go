package trial

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/trial"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	trialType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/trial"

	"github.com/gin-gonic/gin"
)

// SchemaValidateHandle 创建产品
func SchemaValidateHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req trialType.SchemaValidate
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := trial.NewSchemaValidateLogic(c, svcCtx)
		resp, err := logic.SchemaValidate(&req)
		response.HandleResponse(c, resp, err)
	}
}
