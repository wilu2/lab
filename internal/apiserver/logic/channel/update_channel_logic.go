package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/channel"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
)

type UpdateChannelLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateChannelLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UpdateChannelLogic {
	return UpdateChannelLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UpdateChannel 修改渠道信息
func (l *UpdateChannelLogic) UpdateChannel(req *channel.UpdateChannelReq) (resp channel.Channel, err error) {
	// todo: add your logic here and delete this line
	var (
		user = l.ginCtx.Keys["user"].(auth.UserInfo)
		q    = query.Use(l.svcCtx.Db)
	)

	if user.Role == "view" {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	err = q.Transaction(func(tx *query.Query) error {
		var (
			cT = query.Use(l.svcCtx.Db).Channel
		)

		if err := verifyUpdateNameUnique(req.Name, int32(req.ID), l.svcCtx, l.ctx); err != nil {
			return err
		}

		var updateInfo = make(map[string]interface{})
		updateInfo["Name"] = &req.Name
		updateInfo["LastEditorID"] = user.ID
		updateInfo["LastUpdateTime"] = time.Now().In(time.UTC)

		if _, err = cT.WithContext(l.ctx).Where(cT.ID.Eq(int32(req.ID))).Updates(updateInfo); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return
	}

	// itemInfo := model.Service(resInfo)
	resp.ID = req.ID
	resp.Name = req.ChannelDef.Name

	req_body, _ := json.Marshal(req)

	logInfo := optlogs.OptLogger{
		Operation:    "update",
		Resource:     fmt.Sprintf(`%s：%d`, "渠道ID", resp.ID),
		ResourceType: "channel",
		UserID:       user.ID,
		ReqBody:      string(req_body),
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)
	return
}
