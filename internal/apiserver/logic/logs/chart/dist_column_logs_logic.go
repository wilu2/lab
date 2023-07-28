package logs

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/graph"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"
)

type DistColumnLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewDistColumnLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) DistColumnLogic {
	return DistColumnLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// DistColumn 状态码分布柱状图
func (l *DistColumnLogic) DistColumn(req *logs.GraphReq) (resp graph.ColumnChartInfo, err error) {

	var (
		beginDate   = req.DateRange.BeginDate
		endDate     = req.DateRange.EndDate
		user        = l.ginCtx.Keys["user"].(auth.UserInfo)
		category    = make([]string, 0)
		appList     []auth.AppLiteInfo
		routeIDList []string
	)
	appList, routeIDList, err = user.GetAppList(l.svcCtx, int32(req.ServiceID), int32(req.ChannelID), req.ApplicationID, req.BeginDate, req.EndDate)
	if err != nil {
		err = code.WithCodeMsg(code.SumRequestErr)
		return
	}
	allStatus, allData := getAllDist(routeIDList, beginDate, endDate, l.svcCtx, l.ctx)

	for _, status := range allStatus {
		sta := strconv.FormatInt(int64(status), 10)
		category = append(category, sta)
	}

	series := make([]graph.ColumnItem, 0)

	for _, app := range appList {
		series = append(series, graph.ColumnItem{
			Name: app.Name,
			Data: getDist(app.RouteID, allStatus, allData),
		})
	}

	resp.Category = category
	resp.Series = series

	return
}

func getDist(routeID string, allStatus []int32, allData []routeDist) []int {
	data := make([]int, len(allStatus))
	for i, status := range allStatus {
		for _, item := range allData {
			if item.RouteID == routeID && item.Status == status {
				if item.Count != 0 {
					data[i] = int(item.Count)
				}
			}
		}
	}
	return data
}

type routeDist struct {
	RouteID string `gorm:"column:route_id;type:character varying" json:"route_id"`
	Status  int32  `gorm:"column:status;type:integer" json:"status"`
	Count   int32  `gorm:"column:count;type:integer" json:"count"`
}

func getAllDist(application_id []string, beginDate, endDate int64, svcCtx *svc.ServiceContext, ctx context.Context) ([]int32, []routeDist) {
	var (
		lT     = query.Use(svcCtx.Db).AccessLog
		status []int32
		result []routeDist
	)

	_ = lT.WithContext(ctx).
		Where(lT.RouteID.In(application_id...)).
		Where(lT.Timestamp.Between(beginDate, endDate)).
		Order(lT.Status).
		Distinct(lT.Status).
		Pluck(lT.Status, &status)

	_ = lT.WithContext(ctx).
		Select(
			lT.RouteID,
			lT.Status,
			lT.RequestID.Count().As("count"),
		).
		Where(lT.RouteID.In(application_id...)).
		Where(lT.Timestamp.Between(beginDate, endDate)).
		Group(lT.RouteID, lT.Status).
		Order(lT.RouteID).
		Scan(&result)

	return status, result
}
