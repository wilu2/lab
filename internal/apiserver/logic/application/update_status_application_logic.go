package application

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/application"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
)

type UpdateStatusApplicationLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateStatusApplicationLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UpdateStatusApplicationLogic {
	return UpdateStatusApplicationLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UpdateStatusApplication 修改应用信息
func (l *UpdateStatusApplicationLogic) UpdateStatusApplication(req *application.UpdateStatusInfo) (resp application.Application, err error) {
	var (
		aT        = query.Use(l.svcCtx.Db).Application
		user      = l.ginCtx.Keys["user"].(auth.UserInfo)
		routeInfo apisix.Route
	)

	if user.Role == `view` {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	application, _ := aT.WithContext(l.ctx).Where(aT.ID.Eq(int32(req.ID))).First()

	err = apisix.UpdateRouteStatusDb(l.svcCtx.Db, application.RouteID, req.Status)
	if err != nil {
		return
	}
	aT.WithContext(l.ctx).Where(aT.ID.Eq(int32(req.ID))).UpdateSimple(aT.Status.Value(req.Status))

	resp.RouteDef = routeInfo

	var operation string

	if req.Status == 0 {
		operation = "offline"
	} else {
		operation = "online"
	}

	req_body, _ := json.Marshal(req)

	logInfo := optlogs.OptLogger{
		Operation:    operation,
		Resource:     fmt.Sprintf(`%s：%d`, "应用ID", req.ID),
		ResourceType: "application",
		UserID:       user.ID,
		ReqBody:      string(req_body),
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)
	return
}
