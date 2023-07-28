package body_logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"
)

type GetResponseBodyLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewGetResponseBodyLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) GetResponseBodyLogic {
	return GetResponseBodyLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// GetResponseBody 获取返回值
func (l *GetResponseBodyLogic) GetResponseBody(req *logs.BodyLoggerReq) (resp logs.RespBody, err error) {
	var (
		lT = query.Use(l.svcCtx.Db).AccessLog
	)

	query := lT.WithContext(l.ctx).
		Select(
			lT.Timestamp,
		).Where(lT.RequestID.Eq(req.RequestID))

	result, err := query.First()

	if err != nil {
		err = code.WithCodeMsg(code.UserLoginErr)
		return
	}

	date := time.Unix(result.Timestamp, 0).Format("2006-01-02")

	content, rErr := os.ReadFile(fmt.Sprintf("%s/%s/%s", viper.GetString("control.body-logger-dir"), date, req.RequestID+`.resp`))
	if rErr != nil {
		return
	}
	resp.RespBody = string(content)
	return
}
