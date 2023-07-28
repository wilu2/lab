package channel

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/channel"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	channelType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/channel"

	"github.com/gin-gonic/gin"
)

// UpdateChannelHandle 修改渠道信息
func UpdateChannelHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req channelType.UpdateChannelReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := channel.NewUpdateChannelLogic(c, svcCtx)
		resp, err := logic.UpdateChannel(&req)
		response.HandleResponse(c, resp, err)
	}
}
