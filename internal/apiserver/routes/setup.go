package routes

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/routes/application"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/routes/auth"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/routes/channel"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/routes/config"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/routes/document"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/routes/logs"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/routes/opt_logs"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/routes/service"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/routes/summary"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/routes/trial"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/routes/user"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

func Setup(e *gin.Engine, svcCtx *svc.ServiceContext) {

	summary.RegisterSummaryRoute(e, svcCtx)

	service.RegisterProductRoute(e, svcCtx)

	channel.RegisterChannelRoute(e, svcCtx)

	application.RegisterApplicationRoute(e, svcCtx)

	user.RegisterUserRoute(e, svcCtx)

	auth.RegisterAuthRoute(e, svcCtx)

	trial.RegisterTrialRoute(e, svcCtx)

	opt_logs.RegisterOptLogsRoute(e, svcCtx)

	logs.RegisterLogsRoute(e, svcCtx)

	config.RegisterConfigRoute(e, svcCtx)

	config.RegisterConfigRouteWithoutAuth(e, svcCtx)

	document.RegisterDocumentRoute(e, svcCtx)
}
