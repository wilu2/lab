package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/user"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
)

type UpdateUserLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UpdateUserLogic {
	return UpdateUserLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UpdateUser 修改用户信息
func (l *UpdateUserLogic) UpdateUser(req *user.UpdateUser) (resp user.User, err error) {
	// todo: add your logic here and delete this line
	var (
		user = l.ginCtx.Keys["user"].(auth.UserInfo)
		q    = query.Use(l.svcCtx.Db)
	)

	if user.Role == "view" {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	if user.ID == req.ID {
		err = code.WithCodeMsg(code.CantEditYourself)
		return
	}

	err = q.Transaction(func(tx *query.Query) error {
		var (
			uT         = query.Use(l.svcCtx.Db).User
			updateInfo = make(map[string]interface{})
			error      error
		)

		if error = verifyUpdateAccountUniq(req.Account, req.ID, l.svcCtx, l.ctx); error != nil {
			return error
		}

		if error = verifyChannel(req.Channels, req.ID, l.svcCtx, l.ctx); error != nil {
			return error
		}

		targetInfo, error := uT.WithContext(l.ctx).Where(uT.ID.Eq(req.ID)).Find()
		if error != nil {
			return error
		}

		if targetInfo[0].Role == `admin` {
			acQuery := uT.WithContext(l.ctx).Where(uT.Role.Eq("admin"), uT.Abandoned.Is(false))
			aCount, _ := acQuery.Count()
			if aCount < 2 && ((req.UserDef.Status != nil && !*req.UserDef.Status) || req.UserDef.Role != "admin") {
				error = code.WithCodeMsg(code.OnlyLastAdmin)
				return error
			}
		}

		if user.Role == "admin" {
			if req.UserDef.Account != "" {
				updateInfo["Account"] = &req.UserDef.Account
			}

			if req.UserDef.Password != "" {
				updateInfo["Password"] = &req.UserDef.Password
			}

			if req.UserDef.Alias != "" {
				updateInfo["Alias_"] = &req.UserDef.Alias
			}
			if req.UserDef.Desc != "" {
				updateInfo["Description"] = &req.UserDef.Desc
			}

			if req.UserDef.Role != "" {
				updateInfo["Role"] = &req.UserDef.Role
			}
			if req.UserDef.Status != nil {
				updateInfo["Status"] = &req.UserDef.Status
			}
		}

		var channels []byte
		if channels, error = json.Marshal(req.Channels); error != nil {
			return error
		}

		if channels != nil {
			updateInfo["Channels"] = channels
		}

		if _, error = uT.WithContext(l.ctx).Where(uT.ID.Eq(int32(req.ID))).Updates(updateInfo); error != nil {
			return error
		}

		resp.ID = req.ID
		resp.Account = req.UserDef.Account

		return nil
	})

	if err != nil {
		return
	}

	req.Password = ""
	req_body, _ := json.Marshal(req)

	logInfo := optlogs.OptLogger{
		Operation:    "update",
		Resource:     fmt.Sprintf(`%s：%d`, "用户ID", resp.ID),
		ResourceType: "user",
		UserID:       user.ID,
		ReqBody:      string(req_body),
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)
	return
}
