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
)

type UpdateApplicationLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateApplicationLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UpdateApplicationLogic {
	return UpdateApplicationLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UpdateApplication 修改应用信息
func (l *UpdateApplicationLogic) UpdateApplication(req *application.UpdateApplicationInfo) (resp application.Application, err error) {
	// todo: add your logic here and delete this line
	var (
		user     = l.ginCtx.Keys["user"].(auth.UserInfo)
		q        = query.Use(l.svcCtx.Db)
		routeDef = req.RouteDef
	)

	if user.Role == `view` {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	err = q.Transaction(func(tx *query.Query) error {
		var (
			aT                 = query.Use(l.svcCtx.Db).Application
			sT                 = query.Use(l.svcCtx.Db).Service
			oriRouteInfo       = apisix.Route{}
			genRouteDef        = req.RouteDef
			oriRouteID         string
			oriApplicationInfo *model.Application
			serviceInfo        *model.Service
		)

		// 1.获取应用配置信息中的路由ID
		if oriApplicationInfo, err = aT.WithContext(l.ctx).
			Select(aT.RouteID).
			Where(aT.ID.Eq(int32(req.ID))).First(); err != nil {
			return err
		} else if oriApplicationInfo.RouteID == "" {
			return fmt.Errorf("route id is nil")
		} else {
			oriRouteID = oriApplicationInfo.RouteID
		}

		oriRouteInfo, err := apisix.GetRouteDb(l.svcCtx.Db, oriRouteID)
		if err != nil {
			return err
		}

		//2.验证应用名的唯一性
		if err = verifyUpdateNameUnique(routeDef.Name, int32(req.ID), l.svcCtx, l.ctx); err != nil {
			return err
		}

		if serviceInfo, err = sT.WithContext(l.ctx).
			Select(sT.UpstreamID).
			Where(sT.ID.Eq(int32(req.ServiceID))).First(); err != nil {
			return err
		}

		//3. 生成新的路由配置
		genRouteDef.UpstreamID = serviceInfo.UpstreamID
		genRouteDef.BaseInfo.ID = oriRouteInfo.ID

		if genRouteDef.Plugins["limit-count"] != nil {
			genRouteDef.Plugins = configLimitCount(routeDef.Plugins)
		}
		if genRouteDef.Plugins["client-control"] != nil {
			genRouteDef.Plugins = configClientControl(routeDef.Plugins)
		}
		genRouteDef.Plugins = enableSqlLogger(routeDef.Plugins)
		genRouteDef.Plugins = enableBodyLogger(routeDef.Plugins)

		_, err = apisix.UpdateRouteDb(l.svcCtx.Db, &genRouteDef)
		if err != nil {
			return err
		}

		//4. 更新数据库
		var updateInfo = make(map[string]interface{})
		if req.ApplicationDef.RouteDef.Name != "" {
			updateInfo["Name"] = &req.ApplicationDef.RouteDef.Name
		}
		if req.ApplicationDef.ChannelID != 0 {
			updateInfo["ChannelID"] = &req.ApplicationDef.ChannelID
		}
		if req.ApplicationDef.ServiceID != 0 {
			updateInfo["ServiceID"] = &req.ApplicationDef.ServiceID
		}
		// if req.ApplicationDef.GroupID != 0 {
		// 	updateInfo["GroupID"] = &req.ApplicationDef.GroupID
		// }
		updateInfo["LastEditorID"] = user.ID
		updateInfo["LastUpdateTime"] = time.Now().In(time.UTC)

		if _, err := aT.WithContext(l.ctx).Where(aT.ID.Eq(int32(req.ID))).Updates(updateInfo); err != nil {
			apisix.UpdateRouteDb(l.svcCtx.Db, &oriRouteInfo)
			return err
		}

		resp.ID = req.ID
		resp.RouteDef = genRouteDef
		resp.ChannelID = req.ApplicationDef.ChannelID
		resp.ServiceID = req.ApplicationDef.ServiceID
		return nil
	})

	if err != nil {
		return
	}

	req_body, _ := json.Marshal(req)

	logInfo := optlogs.OptLogger{
		Operation:    "update",
		Resource:     fmt.Sprintf(`%s：%d`, "应用ID", req.ID),
		ResourceType: "application",
		UserID:       user.ID,
		ReqBody:      string(req_body),
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)

	return
}
