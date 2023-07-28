package application

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/application"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	applicationType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/application"

	"github.com/gin-gonic/gin"
)

// GetApplicationHandle 获取应用信息
func GetApplicationHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req applicationType.IdPathParam
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := application.NewGetApplicationLogic(c, svcCtx)
		resp, err := logic.GetApplication(&req)
		response.HandleResponse(c, resp, err)
	}
}
