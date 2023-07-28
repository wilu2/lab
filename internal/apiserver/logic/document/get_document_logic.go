package document

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal/query"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/document"
)

type GetDocumentLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewGetDocumentLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) GetDocumentLogic {
	return GetDocumentLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// GetApplication 获取应用信息
func (l *GetDocumentLogic) GetDocument(req *document.IdPathParam) (resp string, err error) {
	var (
		dT = query.Use(l.svcCtx.Db).Document
		// u      = ginCtx.Keys["user"].(*model.User)
	)

	item, _ := dT.WithContext(l.ctx).Where(dT.ID.Eq(int32(req.ID))).First()
	if item == nil {
		return resp, code.WithCodeMsg(code.NotFound, "")

	}
	resp = item.Content
	return
}
