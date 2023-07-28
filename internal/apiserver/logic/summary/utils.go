package summary

import (
	"context"

	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
)

type AppInfo struct {
	ID          string `gorm:"column:id;type:character varying" json:"id"`
	Name        string `gorm:"column:name;type:character varying" json:"name"`
	RouteID     string `gorm:"column:route_id;type:character varying" json:"route_id"`
	Abandoned   bool   `gorm:"column:abandoned;type:boolean" json:"abandoned"`
	ChannelName string `gorm:"column:channel_name;type:character varying" json:"channel_name"`
}

func getAppList(beginDate, endDate int64, user auth.UserInfo, svcCtx *svc.ServiceContext, ctx context.Context) ([]AppInfo, []string) {
	var (
		logT      = query.Use(svcCtx.Db).AccessLog
		aT        = query.Use(svcCtx.Db).Application
		cT        = query.Use(svcCtx.Db).Channel
		appDetail []AppInfo
	)

	aQuery := aT.WithContext(ctx).
		Select(
			aT.Name,
			aT.RouteID,
			cT.Name.As("channel_name"),
		).LeftJoin(cT, cT.ID.EqCol(aT.ChannelID))

	if user.Role != "admin" {
		aQuery = aQuery.
			Where(aT.Abandoned.Is(false))
	}

	if user.Role == "view" { // 当前用户下的 channel
		aQuery = aQuery.Where(aT.ChannelID.In(user.Channels...))
	}

	aQuery.Scan(&appDetail)

	routeIDList := make([]string, 0)
	for _, app := range appDetail {
		routeIDList = append(routeIDList, app.RouteID)
	}

	// 筛选 AppID 对应的 routeId 是否有数据
	_ = logT.WithContext(ctx).
		Select(logT.RouteID, logT.RequestID.Count()).
		Where(logT.Timestamp.Between(beginDate, endDate)).
		Where(logT.RouteID.In(routeIDList...)).
		Group(logT.RouteID).
		Having(logT.RequestID.Count().Gt(0)).
		Distinct(logT.RouteID).
		Pluck(logT.RouteID, &routeIDList)

	appList := make([]AppInfo, 0)
	routeIDMap := make(map[string]struct{}, 0)
	for _, routeID := range routeIDList {
		routeIDMap[routeID] = struct{}{}
	}

	for _, app := range appDetail {
		if _, ok := routeIDMap[app.RouteID]; ok {
			appList = append(appList, AppInfo{
				Name:    app.Name + "(" + app.ChannelName + ")",
				RouteID: app.RouteID,
			})
		}
	}

	return appList, routeIDList
}
