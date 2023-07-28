package logs

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/consts"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/graph"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"
	"gitlab.intsig.net/textin-gateway/pkg/log"
	"gorm.io/gen/field"
)

type AvgRespTimeLineLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewAvgRespTimeLineLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) AvgRespTimeLineLogic {
	return AvgRespTimeLineLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

func (l *AvgRespTimeLineLogic) AvgRespTimeLine(req *logs.GraphReq) (resp graph.LineChartInfo, err error) {
	var (
		user        = l.ginCtx.Keys["user"].(auth.UserInfo)
		appT        = query.Use(l.svcCtx.Db).Application
		appList     []auth.AppLiteInfo
		routeIDList []string
	)

	appList, routeIDList, err = user.GetAppList(l.svcCtx, int32(req.ServiceID), int32(req.ChannelID), req.ApplicationID, req.BeginDate, req.EndDate)
	if err != nil {
		err = code.WithCodeMsg(code.SumRequestErr)
		return
	}
	xAxis := GetAxis(req.BeginDate, req.EndDate, req.Kind) // xAxis 是坐标时间 x 轴
	series, err := l.GetAvgSeries(routeIDList, req.BeginDate, req.EndDate, req.Kind, appList, xAxis, false)

	if user.Role == consts.RoleAdmin && req.ServiceID == 0 && req.ChannelID == 0 && len(req.ApplicationID) == 0 {
		appT.WithContext(l.ctx).Select(appT.RouteID).Pluck(appT.RouteID, &routeIDList)
		seriesOther, err := l.GetAvgSeries(routeIDList, req.BeginDate, req.EndDate, req.Kind, appList, xAxis, true)
		if err != nil {
			return graph.LineChartInfo{}, code.WithCodeMsg(code.SumRequestErr)
		}
		series = append(series, seriesOther...)
	}

	resp.XAxis = xAxis
	resp.Series = series
	return
}

type routeAvgTimeCost struct {
	RouteID     string  `gorm:"column:route_id;type:character varying" json:"route_id"`
	AvgTimeCost float64 `gorm:"column:avg_time_cost;type:float" json:"avg_time_cost"`
	Frac        *int64  `gorm:"column:frac;type:bigint" json:"frac"`
}

func (l *AvgRespTimeLineLogic) GetAvgSeries(
	routeIds []string,
	beginDate, endDate int64,
	kind string,
	appList []auth.AppLiteInfo,
	xAxis []int64,
	other bool,
) (series []graph.LineItem, err error) {
	var (
		ctx            = context.Background()
		logT           = query.Use(l.svcCtx.Db).AccessLog
		groupField     = []field.Expr{}
		selectField    = []field.Expr{}
		logData        []routeAvgTimeCost         // 从 access_log 中查询出来的日志数量
		keyMapTimeCost = make(map[string]float64) // {1_1685894400: 0.01}，如果是 other 的话 {1685894400: 0.01}
	)
	if !other {
		groupField = append(groupField, logT.RouteID)
		selectField = append(selectField, logT.RouteID, logT.TimeCost.Avg().As("avg_time_cost"))
	}

	switch kind {
	case "second":
		groupField = append(groupField, logT.Timestamp)
		selectField = append(selectField, logT.Timestamp.As("frac"))
	case "minute": // 这三个当切大于 60 个点的时候就不能选择，也就是最多 60 小时的数据
		groupField = append(groupField, logT.Minutestamp)
		selectField = append(selectField, logT.Minutestamp.As("frac"))
	case "hour":
		groupField = append(groupField, logT.Hourstamp)
		selectField = append(selectField, logT.Hourstamp.As("frac"))
	case "day", "week", "month":
		groupField = append(groupField, logT.Datestamp)
		selectField = append(selectField, logT.Datestamp.As("frac"))
	case "quarter":
		groupField = append(groupField, logT.Weeklystamp)
		selectField = append(selectField, logT.Weeklystamp.As("frac"))
	case "halfYear", "year":
		groupField = append(groupField, logT.Monthstamp)
		selectField = append(selectField, logT.Monthstamp.As("frac"))
	default:
		groupField = append(groupField, logT.Datestamp)
		selectField = append(selectField, logT.Datestamp.As("frac"))
	}

	logQ := logT.WithContext(ctx).Where(logT.Timestamp.Between(beginDate, endDate))
	if other { // 计算 other
		logQ.Where(logT.RouteID.NotIn(routeIds...))
	} else {
		logQ.Where(logT.RouteID.In(routeIds...))
	}
	err = logQ.Group(groupField...).Select(selectField...).Scan(&logData)
	if err != nil {
		log.Errorf("数据库查询失败 %v", err)
		return
	}

	if len(logData) == 0 { // 没有数据提前返回
		return
	}

	for _, item := range logData {
		var key string
		if other {
			key = fmt.Sprintf("%f", *item.Frac)
		} else {
			key = fmt.Sprintf("%s_%f", item.RouteID, *item.Frac)
		}
		keyMapTimeCost[key] = item.AvgTimeCost
	}
	if other { // 计算 other 其它 appID 的数据
		data := make([]float64, 0)
		for i := 0; i < len(xAxis); i++ {
			data = append(data, keyMapTimeCost[fmt.Sprintf("%d", xAxis[i])])
		}
		series = append(series, graph.LineItem{
			Name: "其他",
			Data: data,
		})
	} else {
		for _, app := range appList { // 计算 appID 的数据，以及获取 app 的 name
			data := make([]float64, 0)
			for i := 0; i < len(xAxis); i++ {
				key := fmt.Sprintf("%s_%f", app.RouteID, xAxis[i])
				data = append(data, keyMapTimeCost[key])
			}
			series = append(series, graph.LineItem{
				Name: app.Name,
				Data: data,
			})
		}
	}
	return
}
