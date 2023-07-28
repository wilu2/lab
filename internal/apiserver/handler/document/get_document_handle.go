package document

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/logic/document"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/response"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	documentType "gitlab.intsig.net/textin-gateway/internal/apiserver/types/document"

	"github.com/gin-gonic/gin"
)

// GetDocumentHandle 获取接口文档
func GetDocumentHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req documentType.IdPathParam
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := document.NewGetDocumentLogic(c, svcCtx)
		resp, err := logic.GetDocument(&req)
		response.HandleStringResponse(c, resp, err)
	}
}
