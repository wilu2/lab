package auth

import (
	"context"

	"gitlab.intsig.net/textin-gateway/internal/apiserver/consts"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/pkg/utils"
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

type LoginInfo struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type TokenInfo struct {
	Token string `json:"token"`
}

type UserInfo struct {
	ID       int32   `json:"id"`
	Account  string  `json:"account"`
	Role     string  `json:"role"` //上游设置
	Channels []int32 `json:"channels"`
}

// GetRole 获取用户信息 channel 和 role
func (userInfo *UserInfo) GetRole() error {
	svcCtx := svc.GetSvc() // 获取 DB
	userT := query.Use(svcCtx.Db).User
	userObj, err := userT.WithContext(context.TODO()).
		Where(userT.ID.Eq(userInfo.ID)).
		Select(userT.Role, userT.Channels).
		First()
	if err != nil {
		log.Errorf("查询用户信息失败: %v", err)
		return err
	}
	userInfo.Channels = utils.ConvertChannels(userObj.Channels)
	userInfo.Role = userObj.Role
	return nil
}

type AppLiteInfo struct {
	ID          string `gorm:"column:id;type:character varying" json:"id"`
	Name        string `gorm:"column:name;type:character varying" json:"name"`
	RouteID     string `gorm:"column:route_id;type:character varying" json:"route_id"`
	Abandoned   bool   `gorm:"column:abandoned;type:boolean" json:"abandoned"`
	ChannelName string `gorm:"column:channel_name;type:character varying" json:"channel_name"`
}

// GetAppList 通过筛选条件，从 AccessLog 表中过滤出来有数据的 AppLiteInfo 信息
func (user *UserInfo) GetAppList(svcCtx *svc.ServiceContext, serviceId, channelId int32, applicationId []int32, beginDate, endDate int64) ([]AppLiteInfo, []string, error) {
	var (
		ctx           = context.Background()
		appT          = query.Use(svcCtx.Db).Application
		channelT      = query.Use(svcCtx.Db).Channel
		routeIDList   []string
		routeIDMapApp = make(map[string]AppLiteInfo, 0)
		appList       = make([]AppLiteInfo, 0)
		err           error
	)
	appTQ := appT.WithContext(ctx).
		Select(
			appT.ID,
			appT.Name,
			appT.RouteID,
			appT.Abandoned,
			channelT.Name.As("channel_name"),
		).
		LeftJoin(channelT, channelT.ID.EqCol(appT.ChannelID))

	if user.Role != consts.RoleAdmin {
		appTQ = appTQ.Where(appT.Abandoned.Is(false))
	}

	if user.Role == consts.RoleView {
		subQuery := appT.WithContext(ctx).
			Columns(appT.ChannelID).
			In(channelT.WithContext(ctx).Select(channelT.ID).Where(channelT.Abandoned.Is(false)))

		appTQ = appTQ.Where(appT.ChannelID.In(user.Channels...), subQuery)
	}

	// appid 和其它两个只能选一个
	if len(applicationId) != 0 {
		appTQ = appTQ.Where(appT.ID.In(applicationId...))
	} else {
		if channelId != 0 {
			appTQ = appTQ.Where(appT.ChannelID.Eq(channelId))
		}
		if serviceId != 0 {
			appTQ = appTQ.Where(appT.ServiceID.Eq(serviceId))
		}
	}
	// 通过筛选条件获取 AppID
	err = appTQ.Scan(&appList)
	if err != nil {
		log.Errorf("获取 Application 表数据失败 %v ", err)
		return nil, nil, err
	}
	for _, app := range appList {
		routeIDList = append(routeIDList, app.RouteID)
		routeIDMapApp[app.RouteID] = app
	}

	// 筛选 AppID 对应的 routeId 是否有数据
	// 数据量多了，没必须要先从 access_log 获取到有数据的 routeID 耗费性能
	appList = make([]AppLiteInfo, 0, len(routeIDList))
	for _, routeID := range routeIDList {
		appList = append(appList, routeIDMapApp[routeID])
	}
	return appList, routeIDList, nil
}
