package application

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/application"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	applicationType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/application"

	"github.com/gin-gonic/gin"
)

// UpdateApplicationHandle 修改应用信息
func UpdateStatusApplicationHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req applicationType.UpdateStatusInfo
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := application.NewUpdateStatusApplicationLogic(c, svcCtx)
		resp, err := logic.UpdateStatusApplication(&req)
		response.HandleResponse(c, resp, err)
	}
}
