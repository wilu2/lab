package service

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/service"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	serviceType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/service"

	"github.com/gin-gonic/gin"
)

// DeleteServiceHandle 删除产品
func DeleteServiceHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req serviceType.IdPathParam
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := service.NewDeleteServiceLogic(c, svcCtx)
		err := logic.DeleteService(&req)
		response.HandleResponse(c, nil, err)
	}
}
