package application

import (
	"github.com/creasty/defaults"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/application"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	applicationType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/application"

	"github.com/gin-gonic/gin"
)

// ListApplicationHandle 获取应用列表
func ListApplicationHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req applicationType.ListReq

		req.Status = -1 // avoid default 0

		if err := defaults.Set(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := application.NewListApplicationLogic(c, svcCtx)
		resp, err := logic.ListApplication(&req)
		response.HandleResponse(c, resp, err)
	}
}
