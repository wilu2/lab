package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/pkg/jwtgen"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
	"gitlab.intsig.net/textin-gateway/internal/pkg/utils"
)

type LoginVerifyLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewLoginVerifyLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) LoginVerifyLogic {
	return LoginVerifyLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// LoginVerify 登录验证
func (l *LoginVerifyLogic) LoginVerify(req *auth.LoginInfo) (resp auth.TokenInfo, err error) {
	// todo: add your logic here and delete this line
	var (
		uT       = query.Use(l.svcCtx.Db).User
		userInfo *model.User
	)

	query := uT.WithContext(l.ctx).Where(uT.Abandoned.Is(false))

	userInfo, err = query.Where(uT.Account.Eq(req.Account), uT.Abandoned.Is(false)).First()

	if err != nil || userInfo == nil {
		err = code.WithCodeMsg(code.UserLoginErr)
		return
	}

	if req.Password != userInfo.Password {
		err = code.WithCodeMsg(code.UserLoginErr)
		return
	}

	channels := utils.ConvertChannels(userInfo.Channels)
	resp.Token = jwtgen.GetToken(auth.UserInfo{
		ID:       userInfo.ID,
		Account:  userInfo.Account,
		Role:     userInfo.Role,
		Channels: channels,
	})
	// resp.Token, err = sso.LoginVerify(req.Account, req.Password)

	itemInfo := optlogs.OptLogger{
		Operation:    "login",
		Resource:     "",
		ResourceType: "user",
		UserID:       userInfo.ID,
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, itemInfo)

	return
}
