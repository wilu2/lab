package logs

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/base"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"
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

// ListService 创建渠道
func (l *ListServiceLogic) ListService(req *logs.ServiceListReq) (resp logs.ServiceList, err error) {
	// todo: add your logic here and delete this line
	var (
		sT   = query.Use(l.svcCtx.Db).Service
		user = l.ginCtx.Keys["user"].(auth.UserInfo)
		// query  = query.Use(l.svcCtx.Db)
	)

	query := sT.WithContext(l.ctx).
		Select(
			sT.ID,
			sT.Name,
			sT.UpstreamID,
			sT.Abandoned,
		)

	if user.Role != `admin` {
		query = query.Where(sT.Abandoned.Is(false))
	}

	items, err := query.Order(sT.ID.Desc()).Debug().Find()

	for _, item := range items {
		var serviceName = item.Name
		if item.Abandoned {
			serviceName = fmt.Sprintf(`%s_%d（已删除）`, item.Name, item.ID)
		}

		resp = append(resp, base.BaseInfo{
			ID:   int(item.ID),
			Name: serviceName,
		})
	}
	return
}
