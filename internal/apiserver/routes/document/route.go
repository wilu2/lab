package document

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/handler/document"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/middleware"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

func RegisterDocumentRoute(e *gin.Engine, svcCtx *svc.ServiceContext) {
	g := e.Group("/gateway/v1")
	g.Use(middleware.AuthMiddleware)
	g.GET("/document/:id", document.GetDocumentHandle(svcCtx))
}
