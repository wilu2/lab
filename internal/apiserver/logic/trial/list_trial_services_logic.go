package trial

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/trial"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

type ListTrialServiceLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewListTrialServiceLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) ListTrialServiceLogic {
	return ListTrialServiceLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// ListTrialService 获取服务列表（体验）
func (l *ListTrialServiceLogic) ListTrialService(req *trial.ListReq) (resp trial.TrialServiceList, err error) {
	var (
		serviceT          = query.Use(l.svcCtx.Db).Service
		streamT           = query.Use(l.svcCtx.Db).ApisixUpstream
		serviceObjList    []*model.Service
		count             int64
		upstreamIDList    = []string{}
		upstreamIDMapName = map[string]string{}
	)

	serviceQ := serviceT.WithContext(l.ctx).
		Select(serviceT.ID, serviceT.UpstreamID).
		Where(serviceT.Abandoned.Is(false)).
		Order(serviceT.ID.Desc())

	if req.Name != "" {
		serviceQ = serviceQ.Where(serviceT.Name.Like("%" + req.Name + "%"))
	}
	serviceObjList, count, err = serviceQ.FindByPage(req.PageSize*(req.Page-1), req.PageSize)
	if err != nil {
		log.Errorf("获取数据库数据失败 %v", err)
		err = code.WithCodeMsg(code.GetTrialListErr)
		return
	}

	resp.Count = uint64(count)
	resp.Items = make([]trial.TrialService, 0, count)
	for _, item := range serviceObjList {
		upstreamIDList = append(upstreamIDList, item.UpstreamID)
	}

	do := streamT.WithContext(l.ctx)
	rows, err := do.
		Select(streamT.Content, streamT.StreamID).
		Where(streamT.Type.Eq(1)).
		Where(streamT.StreamID.In(upstreamIDList...)).Rows()

	if err != nil {
		log.Errorf("获取数据库数据失败 %v", err)
		err = code.WithCodeMsg(code.GetTrialListErr)
		return
	}

	// 从数据库中获取 upstream 的名称
	for rows.Next() {
		var streamObj *model.ApisixUpstream
		streamInfo := apisix.Upstream{}
		do.ScanRows(rows, &streamObj)
		if err = json.Unmarshal([]byte(streamObj.Content), &streamInfo); err != nil {
			log.Errorf("解析 apisix upstream 配置信息失败 %v", err.Error())
			err = code.WithCodeMsg(code.StreamConfigErr)
			return
		}
		upstreamIDMapName[streamObj.StreamID] = streamInfo.Name
	}

	for _, item := range serviceObjList {
		resp.Items = append(resp.Items, trial.TrialService{
			ID:   int(item.ID),
			Name: upstreamIDMapName[item.UpstreamID],
		})
	}
	return
}
