package user

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/user"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) ResetPasswordLogic {
	return ResetPasswordLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// ResetPassword 重置密码
func (l *ResetPasswordLogic) ResetPassword(req *user.ResetPasswordReq) (err error) {
	// todo: add your logic here and delete this line
	var (
		uT       = query.Use(l.svcCtx.Db).User
		user     = l.ginCtx.Keys["user"].(auth.UserInfo)
		userInfo *model.User
	)

	query := uT.WithContext(l.ctx).Where(uT.ID.Eq(int32(user.ID)))

	userInfo, err = query.First()

	if err != nil || userInfo == nil {
		err = code.WithCodeMsg(code.UserResetPwErr)
		return
	}

	if req.OldPassword != userInfo.Password {
		err = code.WithCodeMsg(code.UserResetPwErr)
		fmt.Println(err)
		return
	}

	res, _ := query.Update(uT.Password, req.NewPassword)

	fmt.Println(res)

	// 添加操作日志
	logInfo := optlogs.OptLogger{
		Operation:    "reset_password",
		Resource:     fmt.Sprintf(`%s：%d`, "用户ID", user.ID),
		ResourceType: "user",
		UserID:       user.ID,
		ReqBody:      "",
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)
	return
}
