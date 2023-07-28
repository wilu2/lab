package user

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/user"
)

type GetUserLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) GetUserLogic {
	return GetUserLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// GetChannel 获取渠道信息
func (l *GetUserLogic) GetUser(req *user.IdPathParam) (resp user.User, err error) {
	// todo: add your logic here and delete this line
	var (
		uT   = query.Use(l.svcCtx.Db).User
		user = l.ginCtx.Keys["user"].(auth.UserInfo)
	)

	if user.Role == "view" {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	item, _ := uT.WithContext(l.ctx).Where(uT.ID.Eq(int32(req.ID))).First()
	if item == nil {
		return resp, code.WithCodeMsg(code.NotFound, "")

	}

	resp.ID = item.ID
	resp.Account = item.Account

	return
}
