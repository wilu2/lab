package application

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/application"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	applicationType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/application"

	"github.com/gin-gonic/gin"
)

// CreateApplicationHandle 创建应用
func CreateApplicationHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req applicationType.ApplicationDef
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := application.NewCreateApplicationLogic(c, svcCtx)
		resp, err := logic.CreateApplication(&req)
		response.HandleResponse(c, resp, err)
	}
}
