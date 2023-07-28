package user

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/user"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	userType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/user"

	"github.com/gin-gonic/gin"
)

// GetServiceHandle 获取产品信息
func GetUserHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req userType.IdPathParam
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := user.NewGetUserLogic(c, svcCtx)
		resp, err := logic.GetUser(&req)
		response.HandleResponse(c, resp, err)
	}
}
