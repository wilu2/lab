package opt_logs

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/opt_logs"
)

type ListOptLogsLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewListOptLogsLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) ListOptLogsLogic {
	return ListOptLogsLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// ListOptLogs 获取操作日志列表
func (l *ListOptLogsLogic) ListOptLogs(req *opt_logs.ListReq) (resp opt_logs.OptLogs, err error) {
	// todo: add your logic here and delete this line
	var (
		oT      = query.Use(l.svcCtx.Db).Optlog
		uT      = query.Use(l.svcCtx.Db).User
		user    = l.ginCtx.Keys["user"].(auth.UserInfo)
		optLogs = make([]opt_logs.OptLogItem, 0)
		result  []optLogsJoinUsers
	)

	if user.Role != "admin" {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	query := oT.WithContext(l.ctx).
		Select(
			oT.ID,
			oT.OptTime,
			oT.Operation,
			oT.Resource,
			oT.ResourceType,
			oT.UserID,
			uT.Account.As(`user_name`),
		).
		Join(uT, uT.ID.EqCol(oT.UserID)).
		Order(oT.OptTime.Desc())

	if req.UserID != 0 {
		query = query.Where(oT.UserID.Eq(req.UserID))
	}

	if req.Resource != "" {
		query = query.Where(oT.Resource.Like(("%" + req.Resource + "%")))
	}

	count, _ := query.Count()
	query.Scan(&result)

	for _, item := range result {
		optLogs = append(optLogs, opt_logs.OptLogItem{
			ID:           item.ID,
			OptTime:      item.OptTime.In(time.UTC).Unix(),
			Operation:    item.Operation,
			Resource:     item.Resource,
			ResourceType: item.ResourceType,
			UserName:     item.UserName,
		})
	}

	resp.Items = optLogs
	resp.Count = int(count)
	return
}
