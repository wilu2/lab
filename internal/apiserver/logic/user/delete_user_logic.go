package user

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/user"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
)

type DeleteUserLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) DeleteUserLogic {
	return DeleteUserLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// DeleteChannel 删除渠道
func (l *DeleteUserLogic) DeleteUser(req *user.IdPathParam) (err error) {
	// todo: add your logic here and delete this line
	var (
		q    = query.Use(l.svcCtx.Db)
		user = l.ginCtx.Keys["user"].(auth.UserInfo)
		// query  = query.Use(l.svcCtx.Db)
	)

	if user.Role == `view` {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	if user.ID == req.ID {
		err = code.WithCodeMsg(code.CantEditYourself)
		return
	}

	err = q.Transaction(func(tx *query.Query) error {
		var uT = query.Use(l.svcCtx.Db).User
		targetInfo, error := uT.WithContext(l.ctx).Where(uT.ID.Eq(req.ID)).Find()
		if error != nil {
			err = code.WithCodeMsg(code.BadRequest, "")
		}

		if targetInfo[0].Role == `admin` {
			acQuery := uT.WithContext(l.ctx).Where(uT.Role.Eq("admin"), uT.Abandoned.Is(false))
			aCount, _ := acQuery.Count()
			if aCount < 2 {
				error = code.WithCodeMsg(code.OnlyLastAdmin)
				return error
			}
		}

		var updateInfo = make(map[string]interface{})
		updateInfo["Abandoned"] = true
		updateInfo["DelUniqueKey"] = req.ID

		if _, error := uT.WithContext(l.ctx).Where(uT.ID.Eq(int32(req.ID))).Updates(updateInfo); error != nil {
			return error
		}
		return nil
	})

	if err != nil {
		return
	}

	logInfo := optlogs.OptLogger{
		Operation:    "delete",
		Resource:     fmt.Sprintf(`%s：%d`, "用户ID", req.ID),
		ResourceType: "user",
		UserID:       user.ID,
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)

	return
}
