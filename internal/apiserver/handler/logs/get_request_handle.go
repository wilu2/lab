package logs

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/logs/body_logger"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	logsType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"

	"github.com/gin-gonic/gin"
)

// GetRequest 获取请求体
func GetRequestBodyHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {

	return func(c *gin.Context) {

		var req logsType.BodyLoggerReq
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := body_logger.NewGetRequestBodyLogic(c, svcCtx)
		resp, err := logic.GetRequestBody(&req)
		response.HandleBinaryResponse(c, resp.Name, resp.Type, resp.ReqBody, err)
	}
}
