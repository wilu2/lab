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

type ListChannelLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewListChannelLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) ListChannelLogic {
	return ListChannelLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// ListChannel 创建渠道
func (l *ListChannelLogic) ListChannel(req *logs.ChannelListReq) (resp logs.ChannelList, err error) {
	// todo: add your logic here and delete this line
	var (
		cT   = query.Use(l.svcCtx.Db).Channel
		user = l.ginCtx.Keys["user"].(auth.UserInfo)
		// query  = query.Use(l.svcCtx.Db)
	)

	query := cT.WithContext(l.ctx).
		Select(
			cT.ID,
			cT.Name,
			cT.Abandoned,
		)

	if user.Role != "admin" {
		query = query.Where(cT.Abandoned.Is(false))
	}

	if user.Role == "view" {
		query = query.Where(cT.ID.In(user.Channels...))
	}

	items, err := query.Order(cT.ID.Desc()).Debug().Find()

	for _, item := range items {
		var channelName = item.Name
		if item.Abandoned {
			channelName = fmt.Sprintf(`%s_%d（已删除）`, item.Name, item.ID)
		}

		resp = append(resp, base.BaseInfo{
			ID:   int(item.ID),
			Name: channelName,
		})
	}

	return
}
