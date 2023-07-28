package channel

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/channel"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	channelType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/channel"

	"github.com/gin-gonic/gin"
)

// CreateChannelHandle 创建渠道
func CreateChannelHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req channelType.ChannelDef
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := channel.NewCreateChannelLogic(c, svcCtx)
		resp, err := logic.CreateChannel(&req)
		response.HandleResponse(c, resp, err)
	}
}
