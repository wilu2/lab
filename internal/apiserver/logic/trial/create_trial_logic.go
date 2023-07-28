package trial

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/trial"
	"gitlab.intsig.net/textin-gateway/internal/pkg/factory"
	"gitlab.intsig.net/textin-gateway/internal/pkg/utils"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

type CreateTrialLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTrialLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) CreateTrialLogic {
	return CreateTrialLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// CreateTrial 创建体验用例
func (l *CreateTrialLogic) CreateTrial(req *trial.TrialRouteDef) (resp trial.TrialRoute, err error) {
	var (
		serviceT     = query.Use(l.svcCtx.Db).Service
		routeT       = query.Use(l.svcCtx.Db).ApisixRoute
		upstreamInfo apisix.Upstream
		routeInfo    apisix.Route
		routeObj     *model.ApisixRoute
	)
	item, _ := serviceT.WithContext(l.ctx).
		Select(serviceT.APISet, serviceT.UpstreamID).
		Where(serviceT.ID.Eq(int32(req.ServiceID))).
		First()

	if item == nil {
		return resp, code.WithCodeMsg(code.NotFound, "")
	}

	apiSet := utils.ConvertApiSet(item.APISet)
	upstreamInfo, _ = apisix.GetStreamDb(l.svcCtx.Db, item.UpstreamID)
	upstreamID := fmt.Sprintf("%v", upstreamInfo.ID)

	routeObj, err = routeT.WithContext(l.ctx).
		Where(routeT.Content.Like("%" + upstreamInfo.Name + "(Trial)" + "%")).
		Where(routeT.Type.Eq(2)).
		First()

	if err == nil {
		if err = json.Unmarshal([]byte(routeObj.Content), &routeInfo); err != nil {
			log.Errorf("获取 apisix route 配置信息失败 %v", err.Error())
			err = code.WithCodeMsg(code.RouteConfigErr)
			return
		}
	} else {
		if routeInfo, err = apisix.CreateRouteDb(l.svcCtx.Db, factory.GeneTestRoute(upstreamInfo.Name+"(Trial)", upstreamID), 2); err != nil {
			err = code.WithCodeMsg(code.CreateRouteErr)
			return
		}
	}

	appKey, appSecret := utils.ConvertAppKey(routeInfo.Vars)

	resp.ServiceID = req.ServiceID
	resp.ServiceName = upstreamInfo.Name
	resp.ApiSet = apiSet
	resp.AppKey = appKey
	resp.AppSecret = appSecret

	return
}
