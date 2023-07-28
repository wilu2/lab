package user

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/user"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	userType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/user"

	"github.com/gin-gonic/gin"
)

// UpdateUserHandle 修改产品信息
func UpdateUserHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req userType.UpdateUser
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := user.NewUpdateUserLogic(c, svcCtx)
		resp, err := logic.UpdateUser(&req)
		response.HandleResponse(c, resp, err)
	}
}
