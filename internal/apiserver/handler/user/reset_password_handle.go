package user

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/user"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	userType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/user"

	"github.com/gin-gonic/gin"
)

// ResetPasswordHandle 获取产品信息
func ResetPasswordHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req userType.ResetPasswordReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := user.NewResetPasswordLogic(c, svcCtx)
		err := logic.ResetPassword(&req)
		response.HandleResponse(c, nil, err)
	}
}
