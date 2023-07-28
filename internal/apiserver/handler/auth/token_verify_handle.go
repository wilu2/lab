package auth

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	authType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"

	"github.com/gin-gonic/gin"
)

// TokenVerifyHandle Token验证
func TokenVerifyHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req authType.TokenInfo
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := auth.NewTokenVerifyLogic(c, svcCtx)
		resp, err := logic.TokenVerify(&req)
		response.HandleResponse(c, resp, err)
	}
}
