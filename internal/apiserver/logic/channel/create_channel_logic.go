package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/channel"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

type CreateChannelLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewCreateChannelLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) CreateChannelLogic {
	return CreateChannelLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// CreateChannel 创建渠道
func (l *CreateChannelLogic) CreateChannel(req *channel.ChannelDef) (resp channel.Channel, err error) {
	var (
		q        = query.Use(l.svcCtx.Db)
		user     = l.ginCtx.Keys["user"].(auth.UserInfo)
		channelT = query.Use(l.svcCtx.Db).Channel
	)

	if user.Role == "view" {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	if err = verifyCreateNameUnique(req.Name, l.svcCtx, l.ctx); err != nil {
		log.Errorf("创建 channel 失败 %v", err)
		err = code.WithCodeMsg(code.DuplicateChannelName)
		return
	}
	err = q.Transaction(func(tx *query.Query) error {

		channelInfo := model.Channel{
			Name:           req.Name,
			CreatorID:      user.ID,
			Ctime:          time.Now().In(time.UTC),
			LastEditorID:   user.ID,
			LastUpdateTime: time.Now().In(time.UTC),
		}

		if err = channelT.WithContext(l.ctx).Create(&channelInfo); err != nil {
			return err
		}
		resp.ID = int(channelInfo.ID)
		return nil
	})

	if err != nil {
		log.Errorf("创建 channel 失败 %v", err)
		err = code.WithCodeMsg(code.ChannelCreateErr)
		return
	}

	resp.Name = req.Name

	req_body, _ := json.Marshal(req)

	logInfo := optlogs.OptLogger{
		Operation:    "create",
		Resource:     fmt.Sprintf(`%s：%d`, "渠道ID", resp.ID),
		ResourceType: "channel",
		UserID:       user.ID,
		ReqBody:      string(req_body),
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)

	return
}
