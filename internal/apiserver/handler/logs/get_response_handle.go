package logs

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/logs/body_logger"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	logsType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"

	"github.com/gin-gonic/gin"
)

// GetResponse 获取返回值
func GetResponseBodyHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {

	return func(c *gin.Context) {

		var req logsType.BodyLoggerReq
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := body_logger.NewGetResponseBodyLogic(c, svcCtx)
		resp, err := logic.GetResponseBody(&req)
		response.HandleResponse(c, resp, err)
	}
}
