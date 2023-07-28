package config

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/config"
)

type GetConfigLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewGetConfigLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) GetConfigLogic {
	return GetConfigLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// Summary 创建渠道
func (l *GetConfigLogic) GetConfig() (resp config.Config, err error) {
	var (
		cT = query.Use(l.svcCtx.Db).Config
	)

	items, _ := cT.WithContext(l.ctx).Where(cT.ID.Eq(1)).Find()
	if len(items) <= 0 {
		return resp, code.WithCodeMsg(code.NotFound, "")

	}

	resp.Name = *items[0].Name
	resp.LogoData = *items[0].Logo
	resp.Version = viper.GetString("server.version")

	return
}
