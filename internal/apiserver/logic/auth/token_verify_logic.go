package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/pkg/jwtgen"
)

type TokenVerifyLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTokenVerifyLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TokenVerifyLogic {
	return TokenVerifyLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// TokenVerify TOKEN验证
func (l *TokenVerifyLogic) TokenVerify(req *auth.TokenInfo) (resp auth.UserInfo, err error) {
	// todo: add your logic here and delete this line

	tokenJWTInfo, err := jwtgen.ParseToken(req.Token)
	if err != nil {
		return
	}

	resp = tokenJWTInfo.UserInfo
	return
}
