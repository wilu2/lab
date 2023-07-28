package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/service"
	"gitlab.intsig.net/textin-gateway/internal/pkg/optlogs"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
)

type UpdateServiceLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateServiceLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UpdateServiceLogic {
	return UpdateServiceLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UpdateService 修改服务信息
func (l *UpdateServiceLogic) UpdateService(req *service.UpdateService) (resp service.Service, err error) {
	var (
		serviceT     = query.Use(l.svcCtx.Db).Service
		user         = l.ginCtx.Keys["user"].(auth.UserInfo)
		q            = query.Use(l.svcCtx.Db)
		upstreamInfo = apisix.Upstream{}
		serviceInfo  *model.Service
	)

	if user.Role == `view` {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	err = q.Transaction(func(tx *query.Query) error {
		var (
			dT                = query.Use(l.svcCtx.Db).Document
			vT                = query.Use(l.svcCtx.Db).Version
			apiSet            []byte
			versionInfoSet    = make([]*model.Version, 0)
			modVersionInfoSet = make([]*model.Version, 0)
		)
		serviceInfo, err = serviceT.WithContext(l.ctx).
			Where(serviceT.ID.Eq(int32(req.ID))).First()
		if err != nil {
			return err
		}

		if err := verifyUpdateNameUnique(req.UpstreamUpdate.Name, req.ID, l.svcCtx, l.ctx); err != nil {
			return err
		}

		_, err = apisix.UpdateStreamDb(l.svcCtx.Db, req.UpstreamUpdate, serviceInfo.UpstreamID)
		if err != nil {
			return err
		}

		if apiSet, err = json.Marshal(req.ServiceUpdate.ApiSet); err != nil {
			return err
		}

		var updateInfo = make(map[string]interface{})
		updateInfo["Name"] = &req.UpstreamUpdate.Name
		updateInfo["APISet"] = apiSet
		updateInfo["LastEditorID"] = user.ID
		updateInfo["LastUpdateTime"] = time.Now().In(time.UTC)
		updateInfo["ServiceType"] = req.ServiceType

		if _, err = serviceT.WithContext(l.ctx).
			Where(serviceT.ID.Eq(req.ID)).
			Updates(updateInfo); err != nil {
			return err
		}

		if _, err = dT.WithContext(l.ctx).
			Where(dT.ID.Eq(serviceInfo.DocumentID)).
			Update(dT.Content, req.DocumentContent); err != nil {
			return err
		}

		versionQuery := vT.WithContext(l.ctx).
			Select(
				vT.ID,
				vT.ServiceID,
				vT.Version,
				vT.UpstreamMap,
			).Where(vT.ServiceID.Eq(req.ID))

		versionQuery.Scan(&versionInfoSet)

		nodeMap := getNodeMap(versionInfoSet)

		if nodeMapEqual(nodeMap, req.NodeMap) {
			return nil
		}

		if _, err := vT.WithContext(l.ctx).Where(vT.ServiceID.Eq(req.ID)).Delete(); err != nil {
			return err
		}

		versionMap := getVersionMap(req.NodeMap)

		for _, versionInfo := range versionMap {
			var upstreamMap []byte

			if upstreamMap, err = json.Marshal(versionInfo.UpstreamMap); err != nil {
				return err
			}
			modVersionInfoSet = append(modVersionInfoSet, &model.Version{
				ServiceID:   req.ID,
				Version:     versionInfo.Version,
				UpstreamMap: upstreamMap,
			})
		}

		if error := vT.WithContext(l.ctx).Create(modVersionInfoSet...); error != nil {
			return error
		}

		return nil
	})

	if err != nil {
		return
	}

	resp.ID = req.ID
	resp.Upstream = upstreamInfo
	resp.ApiSet = req.ServiceUpdate.ApiSet

	req_body, _ := json.Marshal(req)

	logInfo := optlogs.OptLogger{
		Operation:    "update",
		Resource:     fmt.Sprintf(`%s：%d`, "服务ID", resp.ID),
		ResourceType: "service",
		UserID:       user.ID,
		ReqBody:      string(req_body),
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)

	return
}
