package channel

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/channel"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

type DeleteChannelLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteChannelLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) DeleteChannelLogic {
	return DeleteChannelLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// DeleteChannel 删除渠道
func (l *DeleteChannelLogic) DeleteChannel(req *channel.IdPathParam) (err error) {
	var (
		user     = l.ginCtx.Keys["user"].(auth.UserInfo)
		channelT = query.Use(l.svcCtx.Db).Channel
		appT     = query.Use(l.svcCtx.Db).Application
	)

	if user.Role == "view" {
		return code.WithCodeMsg(code.Forbidden)
	}

	// channel 不存在
	_, err = channelT.WithContext(l.ctx).
		Where(channelT.ID.Eq(req.ID), channelT.Abandoned.Is(false)).
		First()
	if err != nil {
		return code.WithCodeMsg(code.GetChannelErr)
	}

	// 依赖 channel 的 app 不为空
	count, err := appT.WithContext(l.ctx).
		Where(appT.ChannelID.Eq(req.ID), appT.Abandoned.Is(false)).
		Count()
	if err != nil || count > 0 {
		return code.WithCodeMsg(code.ExistRelatedApp)
	}

	// 删除 channel 服务
	_, err = channelT.WithContext(l.ctx).
		Where(channelT.ID.Eq(req.ID)).
		Updates(map[string]interface{}{
			"Abandoned":    true,
			"DelUniqueKey": req.ID,
		})
	if err != nil {
		log.Errorf("删除渠道失败: %v", err)
		return code.WithCodeMsg(code.DeleteChannelErr)
	}

	logInfo := optlogs.OptLogger{
		Operation:    "delete",
		Resource:     fmt.Sprintf(`%s：%d`, "渠道ID", req.ID),
		ResourceType: "channel",
		UserID:       user.ID,
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)

	return
}
