package config

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/config"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

// GetConfigHandle 创建产品
func GetConfigHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		logic := config.NewGetConfigLogic(c, svcCtx)
		resp, err := logic.GetConfig()
		response.HandleResponse(c, resp, err)
	}
}
