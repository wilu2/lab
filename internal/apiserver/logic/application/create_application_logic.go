package application

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
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/application"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

type CreateApplicationLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewCreateApplicationLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) CreateApplicationLogic {
	return CreateApplicationLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// CreateApplication 创建应用
func (l *CreateApplicationLogic) CreateApplication(req *application.ApplicationDef) (resp application.Application, err error) {
	var (
		routeInfo = apisix.Route{}
		user      = l.ginCtx.Keys["user"].(auth.UserInfo)
		q         = query.Use(l.svcCtx.Db)
		aT        = query.Use(l.svcCtx.Db).Application
		sT        = query.Use(l.svcCtx.Db).Service
	)

	if user.Role == `view` {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	req_body, _ := json.Marshal(req)

	if err = verifyCreateNameUnique(req.RouteDef.Name, l.svcCtx, l.ctx); err != nil {
		log.Errorf("创建应用失败 %v", err)
		err = code.WithCodeMsg(code.DuplicateAppName)
		return
	}
	err = q.Transaction(func(tx *query.Query) error {

		routeDef := req.RouteDef
		item, _ := sT.WithContext(l.ctx).Where(sT.ID.Eq(int32(req.ServiceID))).First()

		routeDef.UpstreamID = item.UpstreamID
		routeDef.BaseInfo.ID = apisix.GetFlakeUidStr()

		if routeDef.Plugins["limit-count"] != nil {
			routeDef.Plugins = configLimitCount(routeDef.Plugins)
		}
		routeDef.Plugins = enableSqlLogger(routeDef.Plugins)
		routeDef.Plugins = enableBodyLogger(routeDef.Plugins)

		if routeInfo, err = apisix.CreateRouteDb(l.svcCtx.Db, &routeDef, 1); err != nil {
			return err
		}

		routeID := fmt.Sprintf("%v", routeInfo.ID)
		channelID := int32(req.ChannelID)
		serviceID := int32(req.ServiceID)
		status := int32(routeDef.Status)

		itemInfo := model.Application{
			ChannelID:      channelID,
			ServiceID:      serviceID,
			RouteID:        routeID,
			Name:           routeDef.Name,
			Status:         status,
			CreatorID:      &user.ID,
			Ctime:          time.Now().In(time.UTC),
			LastEditorID:   &user.ID,
			LastUpdateTime: time.Now().In(time.UTC),
		}
		if err = aT.WithContext(l.ctx).Create(&itemInfo); err != nil {
			apisix.DeleteRouteDb(l.svcCtx.Db, routeID)
			return err
		}
		resp.ID = int(itemInfo.ID)
		resp.CreateTime = itemInfo.Ctime.Unix()
		resp.UpdateTime = itemInfo.Ctime.Unix()
		resp.RouteDef = routeInfo
		return nil
	})

	if err != nil {
		log.Errorf("创建应用失败 %v", err)
		err = code.WithCodeMsg(code.AppCreateErr)
		return
	}

	logInfo := optlogs.OptLogger{
		Operation:    "create",
		Resource:     fmt.Sprintf(`%s：%d`, "应用ID", resp.ID),
		ResourceType: "application",
		UserID:       user.ID,
		ReqBody:      string(req_body),
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)

	return
}
