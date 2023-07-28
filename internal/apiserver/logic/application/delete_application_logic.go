package application

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/application"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

type DeleteApplicationLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteApplicationLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) DeleteApplicationLogic {
	return DeleteApplicationLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// DeleteApplication 删除应用
func (l *DeleteApplicationLogic) DeleteApplication(req *application.IdPathParam) (err error) {
	var (
		q      = query.Use(l.svcCtx.Db)
		appT   = query.Use(l.svcCtx.Db).Application
		routeT = query.Use(l.svcCtx.Db).ApisixRoute
		user   = l.ginCtx.Keys["user"].(auth.UserInfo)
		appObj *model.Application
	)

	if user.Role == `view` {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	appObj, err = appT.WithContext(l.ctx).Select(appT.RouteID).Where(appT.ID.Eq(int32(req.ID))).First()
	if err != nil || appObj.Abandoned {
		return code.WithCodeMsg(code.GetAppErr)
	}

	err = q.Transaction(func(tx *query.Query) error {
		// 删除应用对应的 Route
		_, err = routeT.WithContext(l.ctx).Where(routeT.RouteID.Eq(appObj.RouteID)).Delete()
		if err != nil {
			return err
		}

		// 标记 APP 为删除状态
		_, err = appT.WithContext(l.ctx).
			Where(appT.ID.Eq(int32(req.ID))).
			Updates(map[string]interface{}{
				"Abandoned":    true,
				"DelUniqueKey": req.ID,
			})
		return err
	})

	if err != nil {
		log.Errorf("删除应用失败: %v", err)
		err = code.WithCodeMsg(code.DeleteAppErr)
		return
	}
	apisix.GenConfig(l.svcCtx.Db) // 生成 apisix 配置文件

	logInfo := optlogs.OptLogger{
		Operation:    "delete",
		Resource:     fmt.Sprintf(`%s：%d`, "应用ID", req.ID),
		ResourceType: "application",
		UserID:       user.ID,
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)

	return
}
