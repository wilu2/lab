package user

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/user"
	"gitlab.intsig.net/textin-gateway/internal/pkg/utils"
)

type ListUserLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewListUserLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) ListUserLogic {
	return ListUserLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// ListChannel 获取渠道列表
func (l *ListUserLogic) ListUser(req *user.ListReq) (resp user.UserList, err error) {
	// todo: add your logic here and delete this line
	var (
		uT       = query.Use(l.svcCtx.Db).User
		userInfo = l.ginCtx.Keys["user"].(auth.UserInfo)
		// query  = query.Use(l.svcCtx.Db)
	)
	if userInfo.Role == "view" {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	query := uT.WithContext(l.ctx).
		Select(
			uT.ID,
			uT.Account,
			uT.Alias_,
			uT.Role,
			uT.Channels,
			uT.Description,
			uT.Ctime,
			uT.LastUpdateTime,
		).Where(uT.Abandoned.Is(false))

	if req.ID != 0 {
		query = query.Where(uT.ID.Eq(int32(req.ID)))
	}

	if req.Alias != "" {
		query = query.Where(uT.Alias_.Like("%" + req.Alias + "%"))
	}

	if req.Role != "" {
		query = query.Where(uT.Role.Eq(req.Role))
	}

	if req.Account != "" {
		query = query.Where(uT.Account.Like("%" + req.Account + "%"))
	}

	count, _ := query.Count()
	list := query.Limit(req.PageSize).Offset(req.PageSize * (req.Page - 1)).Order(uT.ID)

	items, err := list.Debug().Find()
	resp.Count = uint64(count)
	resp.Items = make([]user.User, 0, len(items))

	for _, item := range items {
		channels := utils.ConvertChannels(item.Channels)

		resp.Items = append(resp.Items, user.User{
			EditInfo: user.EditInfo{
				ID:         item.ID,
				CreateTime: item.Ctime.Unix(),
				UpdateTime: item.LastUpdateTime.Unix(),
			},
			Account:  item.Account,
			Alias:    item.Alias_,
			Role:     item.Role,
			Channels: getChannelList(channels, l.svcCtx, l.ctx),
			Desc:     item.Description,
		})
	}
	return
}
