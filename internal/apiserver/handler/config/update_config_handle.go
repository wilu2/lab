package config

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/config"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	configType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/config"

	"github.com/gin-gonic/gin"
)

// UpdateConfigHandle创建产品
func UpdateConfigHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req configType.UpdateConfig
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := config.NewUpdateConfigLogic(c, svcCtx)
		resp, err := logic.UpdateConfig(&req)
		response.HandleResponse(c, resp, err)
	}
}
