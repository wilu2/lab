package config

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/config"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
)

type UpdateConfigLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateConfigLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UpdateConfigLogic {
	return UpdateConfigLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// Summary 创建渠道
func (l *UpdateConfigLogic) UpdateConfig(req *config.UpdateConfig) (resp config.Config, err error) {
	var (
		cT   = query.Use(l.svcCtx.Db).Config
		user = l.ginCtx.Keys["user"].(auth.UserInfo)
	)

	if user.Role != `admin` {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	var updateInfo = make(map[string]interface{})

	if req.Name != "" {
		updateInfo["Name"] = &req.Name
	}

	if req.LogoData != "" {
		updateInfo["Logo"] = &req.LogoData
	}

	_, err = cT.WithContext(l.ctx).Where(cT.ID.Eq(1)).Updates(updateInfo)
	if err != nil {
		err = code.WithCodeMsg(code.ConfigIconErr, "文件图标过大，请重新上传")
		return
	}

	req_body, _ := json.Marshal(req)

	logInfo := optlogs.OptLogger{
		Operation:    "update",
		Resource:     "",
		ResourceType: "config",
		UserID:       user.ID,
		ReqBody:      string(req_body),
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)
	return
}
