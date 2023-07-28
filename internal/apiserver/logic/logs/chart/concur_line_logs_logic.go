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
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type ConcurLineLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewConcurLineLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) ConcurLineLogic {
	return ConcurLineLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// ConcurLine 并发量折线图
func (l *ConcurLineLogic) ConcurLine(req *logs.GraphReq) (resp graph.LineChartInfo, err error) {
	var (
		user        = l.ginCtx.Keys["user"].(auth.UserInfo)
		appList     []auth.AppLiteInfo
		routeIDList []string
		appT        = query.Use(l.svcCtx.Db).Application
	)

	appList, routeIDList, err = user.GetAppList(l.svcCtx, int32(req.ServiceID), int32(req.ChannelID), req.ApplicationID, req.BeginDate, req.EndDate)
	if err != nil {
		err = code.WithCodeMsg(code.SumRequestErr)
		return
	}
	xAxis := GetAxis(req.BeginDate, req.EndDate, req.Kind) // xAxis 是坐标时间 x 轴
	series, err := l.GetConcurSeries(routeIDList, req.BeginDate, req.EndDate, req.Kind, appList, xAxis, false)
	if err != nil {
		return graph.LineChartInfo{}, code.WithCodeMsg(code.SumRequestErr)
	}

	if user.Role == consts.RoleAdmin && req.ServiceID == 0 && req.ChannelID == 0 && len(req.ApplicationID) == 0 {
		appT.WithContext(l.ctx).Select(appT.RouteID).Pluck(appT.RouteID, &routeIDList)
		seriesOther, err := l.GetConcurSeries(routeIDList, req.BeginDate, req.EndDate, req.Kind, appList, xAxis, true)
		if err != nil {
			return graph.LineChartInfo{}, code.WithCodeMsg(code.SumRequestErr)
		}
		series = append(series, seriesOther...)
	}

	resp.XAxis = xAxis
	resp.Series = series
	return
}

type routeConcur struct {
	RouteID string  `gorm:"column:route_id;type:character varying" json:"route_id"`
	Concur  float64 `gorm:"column:concur;type:double" json:"concur"`
	Frac    *int64  `gorm:"column:frac;type:bigint" json:"frac"`
}

func (l *ConcurLineLogic) GetConcurSeries(
	routeIds []string,
	beginDate, endDate int64,
	kind string,
	appList []auth.AppLiteInfo,
	xAxis []int64,
	other bool,
) (series []graph.LineItem, err error) {
	var (
		ctx         = context.Background()
		logT        = query.Use(l.svcCtx.Db).AccessLog
		groupField  = []field.Expr{}
		selectField = []field.Expr{logT.RouteID.Count().As("count")}
		logData     []routeConcur              // 从 access_log 中查询出来的日志数量
		keyMapCount = make(map[string]float64) // {1_1685894400: 100}，如果是 other 的话 {1685894400: 100}
	)
	if !other { // 非其它数值计算
		groupField = append(groupField, logT.RouteID)
		selectField = append(selectField, logT.RouteID)
	}
	divNum := 1

	switch kind {
	case "second": // 这三个当切大于 60 个点的时候就不能选择，也就是最多 60 小时的数据
		groupField = append(groupField, logT.Timestamp)
		selectField = append(selectField, logT.Timestamp.As("time_second"))
	case "minute":
		groupField = append(groupField, logT.Minutestamp)
		selectField = append(selectField, logT.Minutestamp.As("time_second"))
		divNum = 60
	case "hour":
		groupField = append(groupField, logT.Hourstamp)
		selectField = append(selectField, logT.Hourstamp.As("time_second"))
		divNum = 3600
	case "day":
		groupField = append(groupField, logT.Datestamp)
		selectField = append(selectField, logT.Datestamp.As("time_second"))
		divNum = 86400
	case "week":
		groupField = append(groupField, logT.Weeklystamp)
		selectField = append(selectField, logT.Weeklystamp.As("time_second"))
		divNum = 604800
	case "month":
		groupField = append(groupField, logT.Monthstamp)
		selectField = append(selectField, logT.Monthstamp.As("time_second"))
		divNum = 2592000
	case "year":
		groupField = append(groupField, logT.Yearstamp)
		selectField = append(selectField, logT.Yearstamp.As("time_second"))
		divNum = 86400 * 365
	default:
		groupField = append(groupField, logT.Datestamp)
		selectField = append(selectField, logT.Datestamp.As("time_second"))
		divNum = 86400
	}

	logQ := logT.WithContext(ctx).Where(logT.Timestamp.Between(beginDate, endDate))
	if other { // 计算 other
		logQ.Where(logT.RouteID.NotIn(routeIds...))
	} else {
		logQ.Where(logT.RouteID.In(routeIds...))
	}
	subQuery := logQ.Select(selectField...).Group(groupField...)

	routeIDField := field.NewString("u", "route_id") // 新建临时表的字段
	fracField := field.NewInt32("u", "time_second")
	countField := field.NewInt("u", "count")
	// 时间单位分平均并发量计算：每秒的瞬时请求量去除掉 0 的数据，然后求这个小时的平均值
	if other {
		err = gen.Table(subQuery.As("u")).
			Select(fracField.As("frac"), countField.Avg().As("concur")).
			Group(fracField).
			Scan(&logData)
	} else {
		err = gen.Table(subQuery.As("u")).
			Select(routeIDField, fracField.As("frac"), countField.Avg().As("concur")).
			Group(routeIDField, fracField).
			Scan(&logData)
	}

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
			key = fmt.Sprintf("%d", *item.Frac)
		} else {
			key = fmt.Sprintf("%s_%d", item.RouteID, *item.Frac)
		}
		item.Concur = item.Concur / float64(divNum)
		keyMapCount[key] = item.Concur
	}

	if other { // 计算 other 其它 appID 的数据
		data := make([]float64, 0)
		for i := 0; i < len(xAxis); i++ {
			data = append(data, keyMapCount[fmt.Sprintf("%d", xAxis[i])])
		}
		series = append(series, graph.LineItem{
			Name: "其他",
			Data: data,
		})
	} else {
		for _, app := range appList { // 计算 appID 的数据，以及获取 app 的 name
			data := make([]float64, 0)
			for i := 0; i < len(xAxis); i++ {
				key := fmt.Sprintf("%s_%d", app.RouteID, xAxis[i])
				data = append(data, keyMapCount[key])
			}
			series = append(series, graph.LineItem{
				Name: app.Name,
				Data: data,
			})
		}
	}
	return
}
