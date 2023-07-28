package apisix

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/pkg/log"
	"gorm.io/gorm"
)

// 用于操作 apisix 的 upstream 操作

/*
upstream json 格式 to yaml
  nodes:
    "192.168.60.52:8080": 1
*/

// CreateStreamDb 在数据库中保存 apisix 的 upstream 对象，用 json str 存储
func CreateStreamDb(db *gorm.DB, upstreamDef *UpstreamDef) (Upstream, error) {
	var (
		ctx       = context.Background()
		streamT   = query.Use(db).ApisixUpstream
		streamObj *model.ApisixUpstream
	)
	upstream := Upstream{
		BaseInfo: BaseInfo{
			ID: GetFlakeUidStr(),
		},
		UpstreamDef: *upstreamDef,
	}
	context, _ := json.Marshal(upstream)
	contextYaml, err := upstream.ConvertYAML()
	if err != nil {
		log.Errorf("定义 upstream yaml 失败 %v", err.Error())
		return Upstream{}, err
	}
	streamObj = &model.ApisixUpstream{
		StreamID:    fmt.Sprintf("%v", upstream.ID),
		Content:     string(context),
		Type:        1,
		ContentYaml: string(contextYaml),
	}
	err = streamT.WithContext(ctx).Create(streamObj)
	if err != nil {
		log.Errorf("数据库创建 upstream 失败 %v", err.Error())
		return Upstream{}, err
	}
	GenConfig(db)
	return upstream, nil
}

// DeleteStreamDb 删除数据库中的 apisix route 对象
func DeleteStreamDb(db *gorm.DB, streamId string) error {
	var (
		ctx     = context.Background()
		streamT = query.Use(db).ApisixUpstream
	)
	_, err := streamT.WithContext(ctx).Where(streamT.StreamID.Eq(streamId)).Delete()
	GenConfig(db)
	return err
}

// UpdateStreamDb 修改数据库中 apisix upstream 对象
func UpdateStreamDb(db *gorm.DB, upstreamDef UpstreamDef, upstreamId string) (Upstream, error) {
	var (
		ctx       = context.Background()
		upstreamT = query.Use(db).ApisixUpstream
	)
	upstream := Upstream{
		BaseInfo: BaseInfo{
			ID: upstreamId,
		},
		UpstreamDef: upstreamDef,
	}
	context, err := json.Marshal(upstream)
	if err != nil {
		log.Errorf("定义 upstream 失败 %v", err.Error())
		return Upstream{}, err
	}
	contextYaml, err := upstream.ConvertYAML()
	_, err = upstreamT.WithContext(ctx).Where(upstreamT.StreamID.Eq(upstreamId)).
		Updates(model.ApisixUpstream{Content: string(context), ContentYaml: contextYaml})
	if err != nil {
		log.Errorf("数据库更新 route 失败 %v", err.Error())
		return Upstream{}, err
	}
	GenConfig(db)
	return upstream, nil
}

// GetStreamDb 从数据库中获取 apisix 的 upstream 对象
func GetStreamDb(db *gorm.DB, streamId string) (streamInfo Upstream, err error) {
	var (
		ctx       = context.Background()
		streamT   = query.Use(db).ApisixUpstream
		streamObj *model.ApisixUpstream
	)

	streamQ := streamT.WithContext(ctx).
		Select(streamT.Content).
		Where(streamT.Type.Eq(1)).
		Where(streamT.StreamID.Eq(streamId))
	streamObj, err = streamQ.First()
	if err != nil {
		log.Errorf("获取 apisix upstream 失败 %v", err.Error())
		return
	}
	streamInfo = Upstream{}
	if err = json.Unmarshal([]byte(streamObj.Content), &streamInfo); err != nil {
		log.Errorf("获取 apisix upstream 配置信息失败 %v", err.Error())
		return
	}
	return streamInfo, nil
}

// AllStreamYaml 获取数据库存储 yaml 格式的 stream 数据
func AllStreamYaml(db *gorm.DB) (streamList []string, err error) {
	var (
		ctx           = context.Background()
		streamT       = query.Use(db).ApisixUpstream
		streamObjList []*model.ApisixUpstream
	)
	streamObjList, err = streamT.WithContext(ctx).
		Select(streamT.ContentYaml).
		Where(streamT.Type.Eq(1)).
		Order(streamT.ID).Find()
	if err != nil {
		return nil, err
	}
	for _, streamObj := range streamObjList {
		streamList = append(streamList, streamObj.ContentYaml)
	}
	return
}
