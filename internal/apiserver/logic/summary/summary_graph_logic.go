package summary

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/graph"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/summary"
	"gorm.io/gen/field"
)

type SummaryGraphLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewSummaryGraphLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) SummaryGraphLogic {
	return SummaryGraphLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

type routeSumCount struct {
	RouteID string `gorm:"column:route_id;type:character varying" json:"route_id"`
	Count   int    `gorm:"column:count;type:integer" json:"count"`
	Frac    *int64 `gorm:"column:frac;type:bigint" json:"frac"`
}

// SummaryGraph 总览图像
func (l *SummaryGraphLogic) SummaryGraph(req *summary.SummaryReq) (resp summary.SummaryGraph, err error) {
	var (
		user        = l.ginCtx.Keys["user"].(auth.UserInfo)
		logT        = query.Use(l.svcCtx.Db).AccessLog
		logData     []routeSumCount
		groupField  = []field.Expr{logT.RouteID}
		selectField = []field.Expr{logT.RouteID, logT.RouteID.Count().As("count")}
		series      = make([]graph.PieItem, 0)
		xAxis       = make([]int64, 0)
		tBegin      = time.Unix(req.BeginDate, 0)
		baseLine    = tBegin // 时间类型
		tEnd        = time.Unix(req.EndDate, 0)
		endLine     = tEnd
		unit        = time.Hour
	)
	appList, routeIDList := getAppList(req.BeginDate, req.EndDate, user, l.svcCtx, l.ctx)

	switch req.Type { // 构造筛选条件，和 x轴的间距
	case "week", "month":
		groupField = append(groupField, logT.Datestamp)
		selectField = append(selectField, logT.Datestamp.As("frac"))
		baseLine = time.Date(tBegin.Year(), tBegin.Month(), tBegin.Day(), 0, 0, 0, 0, tBegin.Location())
		unit = time.Hour * 24
	case "quarter":
		groupField = append(groupField, logT.Weeklystamp)
		selectField = append(selectField, logT.Weeklystamp.As("frac"))
		// 这是一个左闭右开的时间区间
		baseLine = time.Date(tBegin.Year(), tBegin.Month(), tBegin.Day()-int(tBegin.Weekday())+1, 0, 0, 0, 0, tBegin.Location())
		endLine = time.Date(tEnd.Year(), tEnd.Month(), tEnd.Day()-int(tEnd.Weekday())+8, 0, 0, 0, 0, tEnd.Location())
		if tBegin.Weekday() == time.Sunday { // 本周的第一天 Weekday 是 0
			baseLine = baseLine.AddDate(0, 0, -7)
		}
		if tEnd.Weekday() == time.Sunday {
			endLine = endLine.AddDate(0, 0, -7)
		}
		unit = time.Hour * 24 * 7
	case "halfYear", "year":
		groupField = append(groupField, logT.Monthstamp)
		selectField = append(selectField, logT.Monthstamp.As("frac"))
		baseLine = time.Date(tBegin.Year(), tBegin.Month(), 1, 0, 0, 0, 0, tBegin.Location())
		endLine = time.Date(tEnd.Year(), tEnd.Month(), 1, 0, 0, 0, 0, tEnd.Location())
	default:
		groupField = append(groupField, logT.Datestamp)
		selectField = append(selectField, logT.Datestamp.As("frac"))
		baseLine = time.Date(tBegin.Year(), tBegin.Month(), tBegin.Day(), 0, 0, 0, 0, tBegin.Location())
		unit = time.Hour * 24
	}

	err = logT.WithContext(l.ctx).
		Where(logT.Timestamp.Between(req.BeginDate, req.EndDate)).
		Where(logT.RouteID.In(routeIDList...)).
		Group(groupField...).
		Order(groupField...).
		Select(selectField...).
		Scan(&logData)
	if err != nil {
		err = code.WithCodeMsg(code.SumRequestErr)
		return
	}
	if req.Type == "halfYear" || req.Type == "year" {
		for baseLine.Before(endLine) || baseLine.Equal(endLine) {
			xAxis = append(xAxis, baseLine.Unix())
			baseLine = time.Date(baseLine.Year(), baseLine.Month()+1, 1, 0, 0, 0, 0, baseLine.Location())
		}
	} else {
		for tm := baseLine.Unix(); tm <= endLine.Unix(); tm += int64(unit.Seconds()) {
			xAxis = append(xAxis, tm)
		}
	}

	keyMapCount := make(map[string]int) // {1_1685894400: 100}，如果是 other 的话 {1685894400: 100}
	for _, item := range logData {
		key := fmt.Sprintf("%s_%d", item.RouteID, *item.Frac)
		keyMapCount[key] = item.Count
	}

	for _, app := range appList {
		data := make([]int, 0)
		for i := 0; i < len(xAxis); i++ {
			key := fmt.Sprintf("%s_%d", app.RouteID, xAxis[i])
			data = append(data, keyMapCount[key])
		}
		series = append(series, graph.PieItem{
			Name: app.Name,
			Data: data,
		})
	}

	resp.YAxis = `应用`
	resp.XAxis = xAxis
	resp.Series = series
	return
}
