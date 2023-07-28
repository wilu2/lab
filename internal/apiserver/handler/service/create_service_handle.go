package service

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/service"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	serviceType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/service"

	"github.com/gin-gonic/gin"
)

// CreateServiceHandle 创建产品
func CreateServiceHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req serviceType.ServiceDef
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := service.NewCreateServiceLogic(c, svcCtx)
		resp, err := logic.CreateService(&req)
		response.HandleResponse(c, resp, err)
	}
}
