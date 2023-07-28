package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/service"
	"gitlab.intsig.net/textin-gateway/internal/pkg/utils"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

type ListServiceLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewListServiceLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) ListServiceLogic {
	return ListServiceLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// ListService 获取服务列表
func (l *ListServiceLogic) ListService(req *service.ListReq) (resp service.ServiceList, err error) {
	// todo: add your logic here and delete this line
	var (
		sT   = query.Use(l.svcCtx.Db).Service
		user = l.ginCtx.Keys["user"].(auth.UserInfo)
	)

	if user.Role == "view" {
		err = code.WithCodeMsg(code.Forbidden)
		return
	}

	query := sT.WithContext(l.ctx).
		Select(
			sT.ID,
			sT.Name,
			sT.ServiceType,
			sT.UpstreamID,
			sT.APISet,
			sT.CreatorID,
			sT.Ctime,
			sT.LastEditorID,
			sT.LastUpdateTime,
		).Where(sT.Abandoned.Is(false))

	if req.ID != 0 {
		query = query.Where(sT.ID.Eq(int32(req.ID)))
	}
	if req.Name != "" {
		query = query.Where(sT.Name.Like("%" + req.Name + "%"))
	}

	count, _ := query.Count()
	list := query.Order(sT.ID.Desc())
	if req.PageSize != 0 || req.Page != 0 {
		list = list.Limit(req.PageSize).Offset(req.PageSize * (req.Page - 1))
	}

	items, err := list.Debug().Find()
	resp.Count = uint64(count)
	resp.Items = make([]service.Service, 0, len(items))

	for _, item := range items {
		upstreamInfo, err := apisix.GetStreamDb(l.svcCtx.Db, item.UpstreamID)
		if err != nil {
			log.Errorf("获取数据错误 %v", err)
			continue
		}

		apiSet := utils.ConvertApiSet(item.APISet)

		resp.Items = append(resp.Items, service.Service{
			EditInfo: service.EditInfo{
				ID:         item.ID,
				CreateTime: item.Ctime.In(time.UTC).Unix(),
				UpdateTime: item.LastUpdateTime.In(time.UTC).Unix(),
			},
			ServiceType: item.ServiceType,
			Upstream:    upstreamInfo,
			ApiSet:      apiSet,
		})
	}

	return
}
