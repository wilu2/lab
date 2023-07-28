package channel

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/channel"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	channelType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/channel"

	"github.com/gin-gonic/gin"
)

// GetChannelHandle 获取渠道信息
func GetChannelHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req channelType.IdPathParam
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := channel.NewGetChannelLogic(c, svcCtx)
		resp, err := logic.GetChannel(&req)
		response.HandleResponse(c, resp, err)
	}
}
