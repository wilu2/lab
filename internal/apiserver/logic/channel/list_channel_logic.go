package channel

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/channel"
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

// ListChannel 获取渠道列表
func (l *ListChannelLogic) ListChannel(req *channel.ListReq) (resp channel.ChannelList, err error) {
	// todo: add your logic here and delete this line
	var (
		cT    = query.Use(l.svcCtx.Db).Channel
		user  = l.ginCtx.Keys["user"].(auth.UserInfo)
		items []*model.Channel
		// query  = query.Use(l.svcCtx.Db)
	)

	query := cT.WithContext(l.ctx).
		Select(
			cT.ID,
			cT.Name,
			cT.CreatorID,
			cT.Ctime,
			cT.LastEditorID,
			cT.LastUpdateTime,
		).Where(cT.Abandoned.Is(false))

	if req.ID != 0 {
		query = query.Where(cT.ID.Eq(int32(req.ID)))
	}
	if req.Name != "" {
		query = query.Where(cT.Name.Like("%" + req.Name + "%"))
	}

	if user.Role == "view" {
		query = query.Where(cT.ID.In(user.Channels...))
	}

	count, _ := query.Count()
	query = query.Order(cT.ID.Desc())
	// 这货写代码纯靠运气
	if req.PageSize != 0 || req.Page != 0 {
		query = query.Limit(req.PageSize).Offset(req.PageSize * (req.Page - 1))
	}
	items, err = query.Find()
	if err != nil {
		err = code.WithCodeMsg(code.GetChannelErr)
		return
	}
	resp.Count = uint64(count)
	resp.Items = make([]channel.Channel, 0, len(items))

	for _, item := range items {
		resp.Items = append(resp.Items, channel.Channel{
			EditInfo: channel.EditInfo{
				ID:         int(item.ID),
				CreateTime: item.Ctime.In(time.UTC).Unix(),
				UpdateTime: item.LastUpdateTime.In(time.UTC).Unix(),
			},
			ChannelDef: channel.ChannelDef{
				Name: item.Name,
			},
		})
	}
	return
}
