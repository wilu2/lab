package user

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/user"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	userType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/user"

	"github.com/gin-gonic/gin"
)

// CreateServiceHandle 创建产品
func CreateUserHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req userType.UserDef
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := user.NewCreateUserLogic(c, svcCtx)
		resp, err := logic.CreateUser(&req)
		response.HandleResponse(c, resp, err)
	}
}
