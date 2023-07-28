package logs

import (
	"context"
	"fmt"

	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/pkg/utils"
)

type SerJoinApp struct {
	ID      int32  `gorm:"column:id;type:integer;primaryKey;autoIncrement:true" json:"id"`
	Name    string `gorm:"column:name;type:character varying" json:"name"`
	RouteID string `gorm:"column:route_id;type:character varying" json:"route_id"`
}

func getUpstreamSet(ctx context.Context, svcCtx *svc.ServiceContext, version []string) []string {
	var (
		vT          = query.Use(svcCtx.Db).Version
		upstreamSet []string
	)

	query := vT.WithContext(ctx).
		Select(
			vT.ID,
			vT.Version,
			vT.UpstreamMap,
		).Where(vT.Version.In(version...))

	items, _ := query.Debug().Find()

	for _, item := range items {
		var upstreamMap = utils.ConvertUpstreamMap(item.UpstreamMap)
		for _, upstream := range upstreamMap {
			upstreamSet = append(upstreamSet, fmt.Sprintf(`%s:%d`, upstream.Host, upstream.Port))
		}
	}

	return upstreamSet
}

func getRouteServiceSet(ctx context.Context, svcCtx *svc.ServiceContext) []SerJoinApp {
	var (
		sT       = query.Use(svcCtx.Db).Service
		aT       = query.Use(svcCtx.Db).Application
		routeSer []SerJoinApp
	)

	query := sT.WithContext(ctx).
		Select(
			sT.ID,
			sT.Name,
			aT.RouteID.As("route_id"),
		).LeftJoin(aT, aT.ServiceID.EqCol(sT.ID))

	query.Scan(&routeSer)

	return routeSer
}
