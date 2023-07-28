package application

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/application"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/base"
	"gitlab.intsig.net/textin-gateway/internal/pkg/utils"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
)

type GetApplicationLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewGetApplicationLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) GetApplicationLogic {
	return GetApplicationLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// GetApplication 获取应用信息
func (l *GetApplicationLogic) GetApplication(req *application.IdPathParam) (resp application.ApplicationListItem, err error) {
	// todo: add your logic here and delete this line
	var (
		sT   = query.Use(l.svcCtx.Db).Service
		cT   = query.Use(l.svcCtx.Db).Channel
		aT   = query.Use(l.svcCtx.Db).Application
		item *appJoinSerJoinChn
		// u      = ginCtx.Keys["user"].(*model.User)
	)
	query := aT.WithContext(l.ctx).
		Select(
			aT.ID,
			aT.RouteID,
			aT.ChannelID,
			cT.Name.As("channel_name"),
			aT.ServiceID,
			sT.Name.As("service_name"),
			sT.APISet.As("api_set"),
			sT.DocumentID.As("document_id"),
			aT.CreatorID,
			aT.Ctime,
			aT.LastEditorID,
			aT.LastUpdateTime,
		).Where(aT.Abandoned.Is(false)).
		LeftJoin(sT, sT.ID.EqCol(aT.ServiceID)).
		LeftJoin(cT, cT.ID.EqCol(aT.ChannelID))

	query = query.Where(aT.ID.Eq(int32(req.ID)))

	query.Scan(&item)

	if item == nil {
		return resp, code.WithCodeMsg(code.NotFound, "")

	}
	routeInfo, _ := apisix.GetRouteDb(l.svcCtx.Db, *item.RouteID)

	apiSet := utils.ConvertApiSet(item.ApiSet)

	resp.ID = int(item.ID)
	resp.CreateTime = item.Ctime.In(time.UTC).Unix()
	resp.UpdateTime = item.LastUpdateTime.In(time.UTC).Unix()
	resp.APISet = apiSet
	resp.DocumentID = item.DocumentID
	resp.RouteDef = routeInfo
	resp.ServiceInfo = base.BaseInfo{
		ID:   int(*item.ServiceID),
		Name: *item.ServiceName,
	}
	resp.ChannelInfo = base.BaseInfo{
		ID:   int(*item.ChannelID),
		Name: *item.ChannelName,
	}
	resp.ChannelID = int(*item.ChannelID)
	resp.ServiceID = int(*item.ServiceID)

	return
}
