package service

import (
	"context"
	"time"

	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/service"
	"gitlab.intsig.net/textin-gateway/internal/pkg/utils"
)

type serJoinDoc struct {
	ID              int32      `gorm:"column:id;type:integer;primaryKey;autoIncrement:true" json:"id"`
	Name            *string    `gorm:"column:name;type:character varying" json:"name"`
	UpstreamID      *string    `gorm:"column:upstream_id;type:character varying" json:"upstream_id"`
	ServiceType     *string    `gorm:"column:service_type;type:character varying" json:"service_type"`
	APISet          model.JSON `gorm:"column:api_set;type:jsonb" json:"api_set"`
	DocumentID      *int32     `gorm:"column:document_id;type:integer" json:"document_id"`
	Abandoned       *bool      `gorm:"column:abandoned;type:boolean;not null;default:false" json:"abandoned"`
	DelUniqueKey    *int32     `gorm:"column:del_unique_key;type:integer;not null;default:0" json:"del_unique_key"`
	CreatorID       *int32     `gorm:"column:creator_id;type:integer" json:"creator_id"`
	Ctime           time.Time  `gorm:"column:ctime;type:timestamp without time zone;not null;default:CURRENT_TIMESTAMP" json:"ctime"`
	LastEditorID    *int32     `gorm:"column:last_editor_id;type:integer" json:"last_editor_id"`
	LastUpdateTime  time.Time  `gorm:"column:last_update_time;type:timestamp without time zone;not null;default:CURRENT_TIMESTAMP" json:"last_update_time"`
	DocumentContent *string    `gorm:"column:document_content;type:character varying" json:"document_content"`
}

type versionInfo struct {
	Version     string
	UpstreamMap []service.Upstream
}

func verifyCreateNameUnique(name string, svcCtx *svc.ServiceContext, ctx context.Context) error {
	var (
		sT = query.Use(svcCtx.Db).Service
	)

	count, err := sT.WithContext(ctx).
		Where(sT.Name.Eq(name),
			sT.Abandoned.Is(false)).Count()

	if err != nil {
		return err
	}

	if count > 0 {
		return code.WithCodeMsg(code.DuplicateServiceName)
	}

	return nil
}

func verifyUpdateNameUnique(name string, id int32, svcCtx *svc.ServiceContext, ctx context.Context) error {
	var (
		sT = query.Use(svcCtx.Db).Service
	)

	count, err := sT.WithContext(ctx).
		Where(sT.Name.Eq(name),
			sT.Abandoned.Is(false),
			sT.ID.Neq(id)).Count()

	if err != nil {
		return err
	}

	if count > 0 {
		return code.WithCodeMsg(code.DuplicateServiceName)
	}

	return nil
}

func getVersionMap(nodeMap []service.NodeVersion) (versionMap []versionInfo) {
	var versionSet []string
	for _, node := range nodeMap {
		if utils.StringSliceIncludes(node.Version, versionSet) {
			insertUpstream(node, &versionMap)
		} else {
			insertUpstreamWithVersion(node, &versionMap, &versionSet)
		}
	}
	return
}

func getNodeMap(versionMap []*model.Version) (nodeMap []service.NodeVersion) {
	for _, versionInfo := range versionMap {
		upstreamSet := utils.ConvertUpstreamMap(versionInfo.UpstreamMap)
		for _, upstream := range upstreamSet {
			nodeMap = append(nodeMap, service.NodeVersion{
				Version: versionInfo.Version,
				Host:    upstream.Host,
				Port:    upstream.Port,
			})
		}
	}
	return
}

func insertUpstream(node service.NodeVersion, versionMap *[]versionInfo) {
	for i := range *versionMap {
		if (*versionMap)[i].Version == node.Version {
			(*versionMap)[i].UpstreamMap = append((*versionMap)[i].UpstreamMap, service.Upstream{
				Host: node.Host,
				Port: node.Port,
			})
		}
	}
}

func insertUpstreamWithVersion(node service.NodeVersion, versionMap *[]versionInfo, versionSet *[]string) {
	(*versionMap) = append((*versionMap), versionInfo{
		Version: node.Version,
		UpstreamMap: []service.Upstream{
			{
				Host: node.Host,
				Port: node.Port,
			},
		},
	})
	(*versionSet) = append((*versionSet), node.Version)
}

func nodeMapEqual(a, b []service.NodeVersion) bool {
	if len(a) != len(b) {
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}

	var isEqual = true

	for _, m := range a {
		for _, n := range b {
			if m != n {
				isEqual = false
				break
			}
		}
	}

	for _, m := range b {
		for _, n := range a {
			if m != n {
				isEqual = false
				break
			}
		}
	}

	return isEqual
}
