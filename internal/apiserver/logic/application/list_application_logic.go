package application

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/application"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/base"
	"gitlab.intsig.net/textin-gateway/internal/pkg/utils"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
)

type ListApplicationLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewListApplicationLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) ListApplicationLogic {
	return ListApplicationLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// ListApplication 获取应用列表
func (l *ListApplicationLogic) ListApplication(req *application.ListReq) (resp application.ApplicationList, err error) {
	// todo: add your logic here and delete this line
	var (
		sT = query.Use(l.svcCtx.Db).Service
		cT = query.Use(l.svcCtx.Db).Channel
		aT = query.Use(l.svcCtx.Db).Application
		// query  = query.Use(l.svcCtx.Db)
		items []appJoinSerJoinChn
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
			sT.DocumentID,
			aT.CreatorID,
			aT.Ctime,
			aT.LastEditorID,
			aT.LastUpdateTime,
		).Where(aT.Abandoned.Is(false)).
		LeftJoin(sT, sT.ID.EqCol(aT.ServiceID)).
		LeftJoin(cT, cT.ID.EqCol(aT.ChannelID))

	if req.ServiceID != 0 {
		query = query.Where(aT.ServiceID.Eq(req.ServiceID))
	}

	if req.ChannelID != 0 {
		query = query.Where(aT.ChannelID.Eq(req.ChannelID))
	}

	if req.ID != 0 {
		query = query.Where(aT.ID.Eq(int32(req.ID)))
	}
	if req.Name != "" {
		query = query.Where(aT.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != -1 {
		query = query.Where(aT.Status.Eq(int32(req.Status)))
	}
	if req.ServiceType != "" {
		query = query.Where(sT.ServiceType.Eq(req.ServiceType))
	}

	count, _ := query.Count()
	list := query.Order(aT.ID.Desc())
	if req.PageSize != 0 || req.Page != 0 {
		list = list.Limit(req.PageSize).Offset(req.PageSize * (req.Page - 1))
	}
	list.Scan(&items)
	resp.Count = uint64(count)
	resp.Items = make([]application.ApplicationListItem, 0, len(items))

	for _, item := range items {
		routeInfo, _ := apisix.GetRouteDb(l.svcCtx.Db, *item.RouteID)

		// ser, _ := sT.WithContext(l.ctx).Select(sT.Name, sT.APISet).Where(sT.ID.Eq(*item.ServiceID)).Find()
		apiSet := utils.ConvertApiSet(item.ApiSet)
		resp.Items = append(resp.Items, application.ApplicationListItem{
			Application: application.Application{
				EditInfo: application.EditInfo{
					ID:         int(item.ID),
					CreateTime: item.Ctime.In(time.UTC).Unix(),
					UpdateTime: item.LastUpdateTime.In(time.UTC).Unix(),
				},
				ApplicationDef: application.ApplicationDef{
					RouteDef: routeInfo,
				},
			},
			APISet:     apiSet,
			DocumentID: item.DocumentID,
			ServiceInfo: base.BaseInfo{
				ID:   int(*item.ServiceID),
				Name: *item.ServiceName,
			},
			ChannelInfo: base.BaseInfo{
				ID:   int(*item.ChannelID),
				Name: *item.ChannelName,
			},
		})
	}

	resp.Count = uint64(count)
	return
}
