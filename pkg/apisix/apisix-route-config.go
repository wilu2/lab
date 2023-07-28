package apisix

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/pkg/log"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

// 用于操作 apisix 的 config 操作

// CreateRouteDb 在数据库中保存 apisix 的 route 对象，content 存储 json 格式，content_yaml 存储需要 yaml 支持的格式
func CreateRouteDb(db *gorm.DB, route *Route, routeType int16) (Route, error) {
	var (
		ctx      = context.Background()
		routeT   = query.Use(db).ApisixRoute
		routeObj *model.ApisixRoute
	)
	context, err := json.Marshal(route)
	if err != nil {
		log.Errorf("定义 route 失败 %v", err.Error())
		return Route{}, err
	}
	contextYaml, err := yaml.Marshal(route)
	contextYamlStr := "    " + strings.Replace(string(contextYaml), "\n", "\n    ", -1)

	routeObj = &model.ApisixRoute{
		RouteID:     fmt.Sprintf("%v", route.ID),
		Content:     string(context),
		ContentYaml: contextYamlStr,
		Type:        routeType,
		Status:      1,
	}
	// if routeType == 2 {
	// 	exportTime := time.Time{}
	// 	exportTime = time.Now().Add(600 * time.Second)
	// 	routeObj.ExpireAt = &exportTime // 当路由类型为2时，过期时间为当前时间加上600秒
	// }
	err = routeT.WithContext(ctx).Create(routeObj)
	if err != nil {
		log.Errorf("数据库创建 route 失败 %v", err.Error())
		return Route{}, err
	}
	route.ID = routeObj.RouteID
	GenConfig(db)
	return *route, nil
}

// DeleteRouteDb 删除数据库中的 apisix route 对象
func DeleteRouteDb(db *gorm.DB, routeId string) error {
	var (
		ctx    = context.Background()
		routeT = query.Use(db).ApisixRoute
	)
	_, err := routeT.WithContext(ctx).Where(routeT.RouteID.Eq(routeId)).Delete()
	GenConfig(db)
	return err
}

// UpdateRouteDb 修改数据库中 apisix route 对象
func UpdateRouteDb(db *gorm.DB, route *Route) (Route, error) {
	var (
		ctx    = context.Background()
		routeT = query.Use(db).ApisixRoute
	)
	context, err := json.Marshal(route)
	if err != nil {
		log.Errorf("定义 route 失败 %v", err.Error())
		return Route{}, err
	}
	contextYaml, err := yaml.Marshal(route)
	contextYamlStr := "    " + strings.Replace(string(contextYaml), "\n", "\n    ", -1)
	routeID, ok := route.ID.(string)
	if !ok {
		log.Errorf("获取 route ID 失败 %v", route)
		return Route{}, err
	}
	_, err = routeT.WithContext(ctx).Where(routeT.RouteID.Eq(routeID)).
		Updates(model.ApisixRoute{Content: string(context), ContentYaml: contextYamlStr})
	if err != nil {
		log.Errorf("数据库更新 route 失败 %v", err.Error())
		return Route{}, err
	}
	GenConfig(db)
	return *route, nil
}

// UpdateRouteStatusDb 修改 apisix route 对象的状态
func UpdateRouteStatusDb(db *gorm.DB, routeId string, status int32) error {
	var (
		ctx    = context.Background()
		routeT = query.Use(db).ApisixRoute
	)
	route, err := routeT.WithContext(ctx).Select(routeT.Content).Where(routeT.RouteID.Eq(routeId)).Where(routeT.Type.Eq(1)).First()
	if err != nil {
		log.Errorf("获取 apisix route 失败 %v", err.Error())
		return err
	}
	routeInfo := Route{}
	if err = json.Unmarshal([]byte(route.Content), &routeInfo); err != nil {
		log.Errorf("获取 apisix route 配置信息失败 %v", err.Error())
		return err
	}
	routeInfo.Status = Status(status)
	context, err := json.Marshal(routeInfo)
	if err != nil {
		log.Errorf("定义 route 失败 %v", err.Error())
		return err
	}
	contextYaml, err := yaml.Marshal(routeInfo)
	contextYamlStr := "    " + strings.Replace(string(contextYaml), "\n", "\n    ", -1)
	_, err = routeT.WithContext(ctx).Where(routeT.RouteID.Eq(routeId)).
		Updates(model.ApisixRoute{Content: string(context), ContentYaml: contextYamlStr, Status: int16(status)})
	if err != nil {
		log.Errorf("数据库更新 route 失败 %v", err.Error())
		return err
	}
	GenConfig(db)
	return nil
}

// GetRouteDb 从数据库中获取 apisix 的 route 对象
func GetRouteDb(db *gorm.DB, routeId string) (routeInfo Route, err error) {
	var (
		ctx      = context.Background()
		routeT   = query.Use(db).ApisixRoute
		routeObj *model.ApisixRoute
	)

	routeQ := routeT.WithContext(ctx).
		Select(routeT.Content).
		Where(routeT.Type.Eq(1)).
		Where(routeT.RouteID.Eq(routeId))
	routeObj, err = routeQ.First()
	if err != nil {
		log.Errorf("获取 apisix route 失败 %v", err.Error())
		return
	}
	routeInfo = Route{}
	if err = json.Unmarshal([]byte(routeObj.Content), &routeInfo); err != nil {
		log.Errorf("获取 apisix route 配置信息失败 %v", err.Error())
		return
	}
	return routeInfo, nil
}

// AllRouteYaml 获取数据库存储 yaml 格式的 route 数据
func AllRouteYaml(db *gorm.DB) (routeList []string, err error) {
	var (
		ctx          = context.Background()
		routeT       = query.Use(db).ApisixRoute
		routeObjList []*model.ApisixRoute
	)
	routeObjList, err = routeT.WithContext(ctx).
		Select(routeT.ContentYaml).
		Where(routeT.Status.Eq(1)).
		Order(routeT.ID).Find()
	if err != nil {
		return nil, err
	}
	for _, routeObj := range routeObjList {
		routeList = append(routeList, routeObj.ContentYaml)
	}
	return
}
