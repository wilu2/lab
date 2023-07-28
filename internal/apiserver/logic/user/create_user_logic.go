package user

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
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/user"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
)

type CreateUserLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) CreateUserLogic {
	return CreateUserLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// CreateUser 创建用户
func (l *CreateUserLogic) CreateUser(req *user.UserDef) (resp user.User, err error) {
	// todo: add your logic here and delete this line
	var (
		user     = l.ginCtx.Keys["user"].(auth.UserInfo)
		q        = query.Use(l.svcCtx.Db)
		userInfo model.User
	)

	if user.Role == "view" {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	err = q.Transaction(func(tx *query.Query) error {
		var (
			uT = query.Use(l.svcCtx.Db).User
		)

		if err = verifyCreateAccountUniq(req.Account, l.svcCtx, l.ctx); err != nil {
			return err
		}

		if err = verifyChannel(req.Channels, user.ID, l.svcCtx, l.ctx); err != nil {
			return err
		}

		var channels []byte
		if channels, err = json.Marshal(req.Channels); err != nil {
			return err
		}

		userInfo = model.User{
			Account:        req.Account,
			Password:       req.Password,
			Alias_:         req.Alias,
			Role:           req.Role,
			Channels:       channels,
			Description:    req.Desc,
			Ctime:          time.Now(),
			LastUpdateTime: time.Now(),
		}
		if err = uT.WithContext(l.ctx).Create(&userInfo); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return
	}

	resp.ID = userInfo.ID
	resp.Account = req.Account
	resp.Alias = req.Alias
	resp.Role = req.Role
	resp.Channels = getChannelList(req.Channels, l.svcCtx, l.ctx)
	resp.Desc = req.Desc

	req.Password = ""
	req_body, _ := json.Marshal(req)

	logInfo := optlogs.OptLogger{
		Operation:    "create",
		Resource:     fmt.Sprintf(`%s：%d`, "用户ID", resp.ID),
		ResourceType: "user",
		UserID:       user.ID,
		ReqBody:      string(req_body),
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)

	return
}
