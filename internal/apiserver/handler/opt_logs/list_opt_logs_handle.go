package opt_logs

import (
	"github.com/creasty/defaults"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/opt_logs"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	opt_logsType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/opt_logs"

	"github.com/gin-gonic/gin"
)

// ListOptLogsHandle 获取操作日志列表
func ListOptLogsHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req opt_logsType.ListReq
		if err := defaults.Set(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := opt_logs.NewListOptLogsLogic(c, svcCtx)
		resp, err := logic.ListOptLogs(&req)
		response.HandleResponse(c, resp, err)
	}
}
