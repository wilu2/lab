package service

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/service"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	serviceType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/service"

	"github.com/gin-gonic/gin"
)

// UpdateServiceHandle 修改产品信息
func UpdateServiceHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req serviceType.UpdateService
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := service.NewUpdateServiceLogic(c, svcCtx)
		resp, err := logic.UpdateService(&req)
		response.HandleResponse(c, resp, err)
	}
}
