package trial

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/trial"
)

type SchemaValidateLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaValidateLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) SchemaValidateLogic {
	return SchemaValidateLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// ListTrialService 获取服务列表（体验）
func (l *SchemaValidateLogic) SchemaValidate(req *trial.SchemaValidate) (resp trial.ValidateResult, err error) {
	// todo: add your logic here and delete this line
	resp.Valid = true
	return
}
