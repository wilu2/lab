package opt_logs

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/opt_logs"
)

type GetOptLogsDetailLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewGetOptLogsDetailLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) GetOptLogsDetailLogic {
	return GetOptLogsDetailLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// GetOptLogsDetail 获取操作日志详情
func (l *GetOptLogsDetailLogic) GetOptLogsDetail(req *opt_logs.DetailReq) (resp opt_logs.OptLogsDetail, err error) {
	// todo: add your logic here and delete this line
	var (
		user = l.ginCtx.Keys["user"].(auth.UserInfo)
		oT   = query.Use(l.svcCtx.Db).Optlog
	)

	if user.Role != "admin" {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	query := oT.WithContext(l.ctx).
		Select(oT.ReqBody).
		Where(oT.ID.Eq(req.ID))

	item, _ := query.First()

	if item == nil {
		return resp, code.WithCodeMsg(code.NotFound, "")

	}
	resp.Detail = item.ReqBody
	return
}
