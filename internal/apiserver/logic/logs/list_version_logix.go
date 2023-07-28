package logs

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"
)

type ListVersionLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewListVersionLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) ListVersionLogic {
	return ListVersionLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// ListService 创建渠道
func (l *ListVersionLogic) ListVersion(req *logs.VersionListReq) (resp []string, err error) {
	var (
		vT = query.Use(l.svcCtx.Db).Version
		// user = l.ginCtx.Keys["user"].(auth.UserInfo)
	)

	resp = make([]string, 0)

	query := vT.WithContext(l.ctx).
		Select(
			vT.ID,
			vT.Version,
		).Where(vT.ServiceID.Eq(req.ServiceID))

	items, err := query.Order(vT.ID.Desc()).Debug().Find()

	for _, item := range items {
		resp = append(resp, item.Version)
	}
	return
}
