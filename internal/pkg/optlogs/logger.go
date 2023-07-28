package optlogs

import (
	"context"
	"time"

	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/model"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
)

type OptLogger struct {
	Operation    string
	Resource     string
	ResourceType string
	UserID       int32
	ReqBody      string
}

func AddOptLogs(ctx context.Context, svcCtx *svc.ServiceContext, logger OptLogger) {

	var (
		oT = query.Use(svcCtx.Db).Optlog
	)

	itemInfo := model.Optlog{
		Operation:    logger.Operation,
		OptTime:      time.Now().In(time.UTC),
		Resource:     &logger.Resource,
		ResourceType: &logger.ResourceType,
		UserID:       logger.UserID,
		ReqBody:      logger.ReqBody,
	}
	oT.WithContext(ctx).Create(&itemInfo)
}
