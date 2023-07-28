package logs

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/base"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"
)

type ListApplicationLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewListApplicationLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) ListApplicationLogic {
	return ListApplicationLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// ListApplication 创建渠道
func (l *ListApplicationLogic) ListApplication(req *logs.ApplicationListReq) (resp logs.ApplicationList, err error) {
	// todo: add your logic here and delete this line
	var (
		aT   = query.Use(l.svcCtx.Db).Application
		user = l.ginCtx.Keys["user"].(auth.UserInfo)
		// query  = query.Use(l.svcCtx.Db)
	)

	query := aT.WithContext(l.ctx).
		Select(
			aT.ID,
			aT.RouteID,
			aT.Name,
			aT.Abandoned,
		)
	if user.Role != "admin" {
		query = query.Where(aT.Abandoned.Is(false))
	}

	if user.Role == "view" {
		query = query.Where(aT.ChannelID.In(user.Channels...))
	}

	if req.ServiceID != 0 {
		query = query.Where(aT.ServiceID.Eq(int32(req.ServiceID)))
	}
	if req.ChannelID != 0 {
		query = query.Where(aT.ChannelID.Eq(int32(req.ChannelID)))
	}

	items, err := query.Order(aT.ID.Desc()).Debug().Find()

	for _, item := range items {
		var routeName = item.Name
		if item.Abandoned {
			routeName = fmt.Sprintf(`%s_%d（已删除）`, item.Name, item.ID)
		}

		resp = append(resp, base.BaseInfo{
			ID:   int(item.ID),
			Name: routeName,
		})
	}

	return
}
