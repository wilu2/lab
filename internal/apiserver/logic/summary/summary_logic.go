package summary

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/consts"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/summary"
)

type SummaryLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewSummaryLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) SummaryLogic {
	return SummaryLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

type abandonedCount struct {
	Abandoned bool  `gorm:"column:abandoned;type:boolean" json:"abandoned"`
	Count     int64 `gorm:"column:count;type:integer" json:"count"`
}

// Summary 创建渠道
func (l *SummaryLogic) Summary() (resp summary.Summary, err error) {
	var (
		logT       = query.Use(l.svcCtx.Db).AccessLog
		serviceT   = query.Use(l.svcCtx.Db).Service
		channelT   = query.Use(l.svcCtx.Db).Channel
		appT       = query.Use(l.svcCtx.Db).Application
		user       = l.ginCtx.Keys["user"].(auth.UserInfo)
		sResult    []abandonedCount
		cResult    []abandonedCount
		aResult    []abandonedCount
		routeIDSet = make([]string, 0)
	)

	lqCount := logT.WithContext(l.ctx).Select(logT.RequestID.Count())
	qsCount := serviceT.WithContext(l.ctx).Select(serviceT.ID.Count().As(`count`), serviceT.Abandoned).Group(serviceT.Abandoned)
	qcCount := channelT.WithContext(l.ctx).Select(channelT.ID.Count().As(`count`), channelT.Abandoned).Group(channelT.Abandoned)
	qaCount := appT.WithContext(l.ctx).Select(appT.ID.Count().As(`count`), appT.Abandoned).Group(appT.Abandoned)

	if user.Role == consts.RoleView {
		qcCount = qcCount.Where(channelT.ID.In(user.Channels...))
		qaCount = qaCount.Where(appT.ChannelID.In(user.Channels...))

		_ = appT.WithContext(l.ctx).
			Select(appT.RouteID).
			Where(appT.ChannelID.In(user.Channels...)).
			Pluck(appT.RouteID, &routeIDSet)
		lqCount = lqCount.Where(logT.RouteID.In(routeIDSet...))
	}

	lCount, _ := lqCount.Count()
	qsCount.Scan(&sResult)
	qcCount.Scan(&cResult)
	qaCount.Scan(&aResult)

	resp.RequestCount = lCount
	resp.ServiceCount = getCount(sResult)
	resp.ChannelCount = getCount(cResult)
	resp.ApplicationCount = getCount(aResult)
	return
}

// getCount 区分数据，按照 Abandoned 是否关闭
func getCount(result []abandonedCount) (count summary.DetailCount) {
	for _, r := range result {
		if r.Abandoned {
			count.Stopped = r.Count
		} else {
			count.Current = r.Count
		}
	}
	count.Total = count.Current + count.Stopped
	return
}
