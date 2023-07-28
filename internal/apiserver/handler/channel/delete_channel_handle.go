package channel

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/channel"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	channelType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/channel"

	"github.com/gin-gonic/gin"
)

// DeleteChannelHandle 删除渠道
func DeleteChannelHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req channelType.IdPathParam
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := channel.NewDeleteChannelLogic(c, svcCtx)
		err := logic.DeleteChannel(&req)
		response.HandleResponse(c, nil, err)
	}
}
