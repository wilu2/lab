package auth

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	authType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"

	"github.com/gin-gonic/gin"
)

// LoginVerifyHandle 登录验证
func LoginVerifyHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req authType.LoginInfo
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := auth.NewLoginVerifyLogic(c, svcCtx)
		resp, err := logic.LoginVerify(&req)
		response.HandleResponse(c, resp, err)
	}
}
