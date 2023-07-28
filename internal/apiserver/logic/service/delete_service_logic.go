package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/service"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

type DeleteServiceLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteServiceLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) DeleteServiceLogic {
	return DeleteServiceLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// DeleteService 删除产品
func (l *DeleteServiceLogic) DeleteService(req *service.IdPathParam) (err error) {
	var (
		user       = l.ginCtx.Keys["user"].(auth.UserInfo)
		q          = query.Use(l.svcCtx.Db)
		serviceT   = query.Use(l.svcCtx.Db).Service
		appT       = query.Use(l.svcCtx.Db).Application
		streamT    = query.Use(l.svcCtx.Db).ApisixUpstream
		routeT     = query.Use(l.svcCtx.Db).ApisixRoute
		serviceObj *model.Service
		count      int64
	)

	if user.Role == `view` {
		return code.WithCodeMsg(code.Forbidden)
	}

	serviceObj, err = serviceT.WithContext(l.ctx).
		Select(serviceT.UpstreamID, serviceT.Abandoned).
		Where(serviceT.ID.Eq(int32(req.ID))).
		First()
	if err != nil || serviceObj.Abandoned {
		return code.WithCodeMsg(code.GetServiceErr)
	}

	// 存在依赖 service 的 app 服务
	count, err = appT.WithContext(l.ctx).
		Where(appT.ServiceID.Eq(req.ID), appT.Abandoned.Is(false)).
		Count()
	if err != nil || count > 0 {
		return code.WithCodeMsg(code.ExistRelatedApp)
	}

	err = q.Transaction(func(tx *query.Query) error {
		// 删除 upstream
		_, err = streamT.WithContext(l.ctx).Where(streamT.StreamID.Eq(serviceObj.UpstreamID)).Delete()
		if err != nil {
			return err
		}
		// 删除体验中心的 route
		upstreamStr := fmt.Sprintf(`"upstream_id":"%s"`, serviceObj.UpstreamID)
		_, err = routeT.WithContext(l.ctx).Where(routeT.Content.Like("%" + upstreamStr + "%")).Where(routeT.Type.Eq(2)).Delete()
		if err != nil {
			return err
		}

		// 将 service 服务标记为删除
		_, err = serviceT.WithContext(l.ctx).
			Where(serviceT.ID.Eq(req.ID)).
			Updates(map[string]interface{}{
				"Abandoned":    true,
				"DelUniqueKey": req.ID,
			})
		return err
	})

	if err != nil {
		log.Errorf("删除服务失败: %v", err)
		return code.WithCodeMsg(code.DelServiceErr)
	}
	apisix.GenConfig(l.svcCtx.Db) // 生成 apisix 配置文件

	logInfo := optlogs.OptLogger{
		Operation:    "delete",
		Resource:     fmt.Sprintf(`%s：%d`, "服务ID", req.ID),
		ResourceType: "service",
		UserID:       user.ID,
		ReqBody:      "",
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)

	return
}
