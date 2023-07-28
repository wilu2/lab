package user

import (
	"github.com/creasty/defaults"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/user"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	userType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/user"

	"github.com/gin-gonic/gin"
)

// ListServiceHandle 获取产品列表
func ListUserHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req userType.ListReq
		if err := defaults.Set(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := user.NewListUserLogic(c, svcCtx)
		resp, err := logic.ListUser(&req)
		response.HandleResponse(c, resp, err)
	}
}
