package channel

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/channel"
)

type GetChannelLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewGetChannelLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) GetChannelLogic {
	return GetChannelLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// GetChannel 获取渠道信息
func (l *GetChannelLogic) GetChannel(req *channel.IdPathParam) (resp channel.Channel, err error) {
	var (
		cT   = query.Use(l.svcCtx.Db).Channel
		user = l.ginCtx.Keys["user"].(auth.UserInfo)
	)

	query := cT.WithContext(l.ctx)

	if user.Role == "view" {
		query = query.Where(cT.ID.In(user.Channels...))
	}

	item, _ := query.Where(cT.ID.Eq(int32(req.ID))).First()

	if item == nil {
		return resp, code.WithCodeMsg(code.NotFound, "")

	}

	resp.ID = int(item.ID)
	resp.Name = item.Name

	return
}
