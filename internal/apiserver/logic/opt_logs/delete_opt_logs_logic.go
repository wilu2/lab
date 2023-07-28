package opt_logs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/opt_logs"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
)

type DeleteOptLogsLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOptLogsLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) DeleteOptLogsLogic {
	return DeleteOptLogsLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// DeleteOptLogs 删除操作日志
func (l *DeleteOptLogsLogic) DeleteOptLogs(req *opt_logs.DeleteReq) (err error) {
	var (
		user      = l.ginCtx.Keys["user"].(auth.UserInfo)
		oT        = query.Use(l.svcCtx.Db).Optlog
		beginDate = time.Unix(int64(req.BeginDate), 0).UTC()
		endDate   = time.Unix(int64(req.EndDate), 0).UTC()
	)

	if user.Role != "admin" {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	_, err = oT.WithContext(l.ctx).Where(oT.OptTime.Between(beginDate, endDate)).Delete()
	if err != nil {
		return code.WithCodeMsg(code.DeleteLogsErr)
	}

	req_body, _ := json.Marshal(req)

	logInfo := optlogs.OptLogger{
		Operation:    "delete",
		Resource:     "",
		ResourceType: "opt_logs",
		UserID:       user.ID,
		ReqBody:      string(req_body),
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)
	return
}
