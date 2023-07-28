package service

import (
	"github.com/creasty/defaults"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/service"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	serviceType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/service"

	"github.com/gin-gonic/gin"
)

// ListServiceHandle 获取产品列表
func ListServiceHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req serviceType.ListReq
		if err := defaults.Set(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := service.NewListServiceLogic(c, svcCtx)
		resp, err := logic.ListService(&req)
		response.HandleResponse(c, resp, err)
	}
}
