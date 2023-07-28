package opt_logs

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/opt_logs"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	opt_logsType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/opt_logs"

	"github.com/gin-gonic/gin"
)

// DeleteOptLogsHandle	删除操作日志
func DeleteOptLogsHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req opt_logsType.DeleteReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := opt_logs.NewDeleteOptLogsLogic(c, svcCtx)
		err := logic.DeleteOptLogs(&req)
		response.HandleResponse(c, nil, err)
	}
}
