package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/service"
	"gitlab.intsig.net/textin-gateway/internal/pkg/utils"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
)

type GetServiceLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewGetServiceLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) GetServiceLogic {
	return GetServiceLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// GetService 获取服务信息
func (l *GetServiceLogic) GetService(req *service.IdPathParam) (resp service.Service, err error) {
	var (
		user    = l.ginCtx.Keys["user"].(auth.UserInfo)
		q       = query.Use(l.svcCtx.Db)
		item    *serJoinDoc
		nodeMap []service.NodeVersion
	)

	if user.Role == `view` {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	err = q.Transaction(func(tx *query.Query) error {
		var (
			sT         = query.Use(l.svcCtx.Db).Service
			dT         = query.Use(l.svcCtx.Db).Document
			vT         = query.Use(l.svcCtx.Db).Version
			versionMap []*model.Version
		)

		query := sT.WithContext(l.ctx).
			Select(
				sT.ID,
				sT.UpstreamID,
				sT.APISet,
				sT.DocumentID,
				sT.ServiceType,
				dT.Content.As("document_content"),
				sT.CreatorID,
				sT.Ctime,
				sT.LastEditorID,
				sT.LastUpdateTime,
			).Where(sT.Abandoned.Is(false)).
			LeftJoin(dT, dT.ID.EqCol(sT.DocumentID))

		query = query.Where(sT.ID.Eq(req.ID))

		query.Scan(&item)

		if item == nil {
			return code.WithCodeMsg(code.NotFound, "")
		}

		versionQuery := vT.WithContext(l.ctx).
			Select(
				vT.ID,
				vT.ServiceID,
				vT.Version,
				vT.UpstreamMap,
			).Where(vT.ServiceID.Eq(req.ID))

		versionQuery.Scan(&versionMap)

		nodeMap = getNodeMap(versionMap)

		return nil
	})

	if err != nil {
		return
	}

	upstreamInfo, _ := apisix.GetStreamDb(l.svcCtx.Db, *item.UpstreamID)
	apiSet := utils.ConvertApiSet(item.APISet)

	resp.ID = item.ID
	resp.ServiceType = *item.ServiceType
	resp.CreateTime = item.Ctime.In(time.UTC).Unix()
	resp.UpdateTime = item.LastUpdateTime.In(time.UTC).Unix()
	resp.DocumentID = *item.DocumentID
	resp.DocumentContent = *item.DocumentContent
	resp.Upstream = upstreamInfo
	resp.ApiSet = apiSet
	resp.NodeMap = nodeMap
	return
}
