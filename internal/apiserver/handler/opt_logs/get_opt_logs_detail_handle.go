package opt_logs

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/opt_logs"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	opt_logsType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/opt_logs"

	"github.com/gin-gonic/gin"
)

// GetOptLogsDetail 获取操作日志详情
func GetOptLogsDetailHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req opt_logsType.DetailReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := opt_logs.NewGetOptLogsDetailLogic(c, svcCtx)
		resp, err := logic.GetOptLogsDetail(&req)
		response.HandleResponse(c, resp, err)
	}
}
