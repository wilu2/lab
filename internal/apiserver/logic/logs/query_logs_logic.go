package logs

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/consts"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

type QueryLogsLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewQueryLogsLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) QueryLogsLogic {
	return QueryLogsLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// QueryLogs
func (l *QueryLogsLogic) QueryLogs(req *logs.QueryReq) (resp logs.Logs, err error) {
	var (
		lT         = query.Use(l.svcCtx.Db).AccessLog
		aT         = query.Use(l.svcCtx.Db).Application
		user       = l.ginCtx.Keys["user"].(auth.UserInfo)
		routeIDSet = make([]string, 0)
	)

	query := lT.WithContext(l.ctx).
		Select(
			lT.RouteID,
			lT.RequestID,
			lT.ClientAddr,
			lT.IsoTime,
			lT.Timestamp,
			lT.TimeCost,
			lT.RequestLength,
			lT.Connection,
			lT.ConnectionRequests,
			lT.URI,
			lT.OriRequest,
			lT.QueryString,
			lT.Status,
			lT.BytesSent,
			lT.Referer,
			lT.UserAgent,
			lT.ForwardedFor,
			lT.Host,
			lT.Node,
			lT.Upstream,
		)

	if user.Role == consts.RoleView {
		_ = aT.WithContext(l.ctx).
			Select(aT.RouteID).
			Where(aT.ChannelID.In(user.Channels...)).
			Pluck(aT.RouteID, &routeIDSet)
		query.Where(lT.RouteID.In(routeIDSet...))
	}

	if len(req.ApplicationID) > 0 || req.ServiceID != 0 || req.ChannelID != 0 {

		routeQuery := aT.WithContext(l.ctx).
			Select(aT.RouteID)

		if len(req.ApplicationID) > 0 {
			routeQuery = routeQuery.Where(aT.ID.In(req.ApplicationID...))
		} else {
			if req.ServiceID != 0 {
				routeQuery = routeQuery.Where(aT.ServiceID.Eq(int32(req.ServiceID)))
			}
			if req.ChannelID != 0 {
				routeQuery = routeQuery.Where(aT.ChannelID.Eq(int32(req.ChannelID)))
			}
		}

		var routeList []string
		routeQuery.Pluck(aT.RouteID, &routeList)
		query = query.Where(lT.RouteID.In(routeList...))
	}

	if req.RequestID != "" {
		query = query.Where(lT.RequestID.Eq(req.RequestID))
	}

	if len(req.Status) > 0 {
		query = query.Where(lT.Status.In(req.Status...))
	}

	if len(req.Version) > 0 {
		query = query.Where(lT.Upstream.In(getUpstreamSet(l.ctx, l.svcCtx, req.Version)...))
	}

	query = query.Where(lT.Timestamp.Between(req.BeginDate/1000, req.EndDate/1000)).Order(lT.IsoTime.Desc())

	items, count, err := query.FindByPage(req.PageSize*(req.Page-1), req.PageSize)
	if err != nil {
		log.Errorf("获取数据错误 %v", err)
		err = code.WithCodeMsg(code.LogGetListErr)
		return
	}
	resp.Count = int(count)
	resp.Items = make([]logs.LogItem, 0, len(items))

	for _, item := range items {
		resp.Items = append(resp.Items, logs.LogItem{
			RequestID:          item.RequestID,
			ClientAddr:         item.ClientAddr,
			IsoTime:            item.IsoTime,
			Timestamp:          item.Timestamp,
			TimeCost:           item.TimeCost,
			RequestLength:      item.RequestLength,
			Connection:         item.Connection,
			ConnectionRequests: item.ConnectionRequests,
			Uri:                item.URI,
			OriRequest:         item.OriRequest,
			Status:             item.Status,
			BytesSent:          item.BytesSent,
			Referer:            *item.Referer,
			UserAgent:          item.UserAgent,
			ForwardedFor:       *item.ForwardedFor,
			Host:               *item.Host,
			Node:               *item.Node,
			Upstream:           item.Upstream,
		})
	}
	return
}
