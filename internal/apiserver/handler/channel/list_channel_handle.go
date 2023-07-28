package channel

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/channel"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	channelType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/channel"

	"github.com/creasty/defaults"
	"github.com/gin-gonic/gin"
)

// ListChannelHandle 获取渠道列表
func ListChannelHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req channelType.ListReq
		if err := defaults.Set(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := channel.NewListChannelLogic(c, svcCtx)
		resp, err := logic.ListChannel(&req)
		response.HandleResponse(c, resp, err)
	}
}
