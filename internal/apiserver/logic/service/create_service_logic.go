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
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

type CreateServiceLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewCreateServiceLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) CreateServiceLogic {
	return CreateServiceLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// CreateService 创建产品
func (l *CreateServiceLogic) CreateService(req *service.ServiceDef) (resp service.Service, err error) {
	var (
		q              = query.Use(l.svcCtx.Db)
		upstreamInfo   = apisix.Upstream{}
		user           = l.ginCtx.Keys["user"].(auth.UserInfo)
		sT             = query.Use(l.svcCtx.Db).Service
		dT             = query.Use(l.svcCtx.Db).Document
		vT             = query.Use(l.svcCtx.Db).Version
		versionInfoSet = make([]*model.Version, 0)
		upstreamID     string
		itemInfo       model.Service
	)

	if err = verifyCreateNameUnique(req.UpstreamDef.Name, l.svcCtx, l.ctx); err != nil {
		return
	}

	err = q.Transaction(func(tx *query.Query) error {

		upstreamInfo, err = apisix.CreateStreamDb(l.svcCtx.Db, &req.UpstreamDef)
		if err != nil {
			return err
		}

		documentInfo := model.Document{
			Content: req.DocumentContent,
		}

		if error := dT.WithContext(l.ctx).Create(&documentInfo); error != nil {
			return error
		}
		upstreamID = fmt.Sprintf("%v", upstreamInfo.ID)

		var apiSet []byte
		if apiSet, err = json.Marshal(req.ApiSet); err != nil {
			return err
		}

		itemInfo = model.Service{
			Name:           req.UpstreamDef.Name,
			UpstreamID:     upstreamID,
			ServiceType:    req.ServiceType,
			APISet:         apiSet,
			DocumentID:     documentInfo.ID,
			CreatorID:      &user.ID,
			Ctime:          time.Now().In(time.UTC),
			LastEditorID:   &user.ID,
			LastUpdateTime: time.Now().In(time.UTC),
		}

		if error := sT.WithContext(l.ctx).Create(&itemInfo); error != nil {
			return error
		}

		versionMap := getVersionMap(req.NodeMap)

		for _, versionInfo := range versionMap {
			var upstreamMap []byte

			if upstreamMap, err = json.Marshal(versionInfo.UpstreamMap); err != nil {
				return err
			}
			versionInfoSet = append(versionInfoSet, &model.Version{
				ServiceID:   itemInfo.ID,
				Version:     versionInfo.Version,
				UpstreamMap: upstreamMap,
			})
		}

		if error := vT.WithContext(l.ctx).Create(versionInfoSet...); error != nil {
			return error
		}

		return nil
	})

	if err != nil {
		log.Errorf("创建 service 失败 %v", err)
		err = code.WithCodeMsg(code.ServiceCreateErr)
		return
	}

	resp.ID = itemInfo.ID
	resp.ServiceType = itemInfo.ServiceType
	resp.CreateTime = itemInfo.Ctime.Unix()
	resp.UpdateTime = itemInfo.Ctime.Unix()
	resp.Upstream = upstreamInfo
	resp.ApiSet = req.ApiSet

	req_body, _ := json.Marshal(req)

	logInfo := optlogs.OptLogger{
		Operation:    "create",
		Resource:     fmt.Sprintf(`%s：%d`, "服务ID", resp.ID),
		ResourceType: "service",
		UserID:       user.ID,
		ReqBody:      string(req_body),
	}

	optlogs.AddOptLogs(l.ctx, l.svcCtx, logInfo)
	return
}
