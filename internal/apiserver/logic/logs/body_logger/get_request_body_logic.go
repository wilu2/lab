package body_logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"github.com/spf13/viper"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/logs"
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

type GetRequestBodyLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewGetRequestBodyLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) GetRequestBodyLogic {
	return GetRequestBodyLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// GetRequestBody 获取请求体
func (l *GetRequestBodyLogic) GetRequestBody(req *logs.BodyLoggerReq) (resp logs.ReqBody, err error) {
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

	content, err := os.ReadFile(fmt.Sprintf("%s/%s/%s", viper.GetString("control.body-logger-dir"), date, req.RequestID+`.req`))
	if err != nil {
		log.Errorf("获取图片错误%v", err)
		err = code.WithCodeMsg(code.GetImageErr)
		return
	}

	resp.Name = req.RequestID
	imageType, err := filetype.Get(content)
	if err != nil {
		log.Errorf("获取图片类型错误%v", err)
		err = code.WithCodeMsg(code.GetImageTypeErr)
		return
	}
	resp.Type = imageType.MIME.Value
	resp.ReqBody = content
	return
}
