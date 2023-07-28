package apiserver

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/config"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal"
	"gorm.io/gorm"

	// "gitlab.intsig.net/textin-gateway/internal/apiserver/middleware"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/routes"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/svc"
	genericOptions "gitlab.intsig.net/textin-gateway/internal/pkg/options"
	genericapiserver "gitlab.intsig.net/textin-gateway/internal/pkg/server"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
	"gitlab.intsig.net/textin-gateway/pkg/log"
	"gitlab.intsig.net/textin-gateway/pkg/shutdown"
)

type apiServer struct {
	cfg              *config.Config
	gs               *shutdown.GracefulShutdown
	apisixOption     *genericOptions.ApisixOptions
	dbOption         *genericOptions.DbSQLOptions
	clickhouseOption *genericOptions.ClickHouseOptions
	genericApiServer *genericapiserver.GenericApiServer
}

// createAPIServer 创建 API Server 服务
func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(shutdown.NewPosixSignalManager()) // 优雅关闭系统

	genericConfig, err := buildGenericConfig(cfg) // 创建 GenericAPIServer 配置
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New() // 初始化完成 HTTP 服务
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		cfg:              cfg,
		gs:               gs,
		genericApiServer: genericServer,
		apisixOption:     cfg.ApisixOptions,
		dbOption:         cfg.DbSQLOptions,
		clickhouseOption: cfg.ClickHouseOptions,
	}

	return server, nil
}

// Run  业务代码路由地址
func (s *apiServer) Run() error {
	s.initDbStore()

	svcCtx := svc.NewServiceContext(s.cfg)
	routes.Setup(s.genericApiServer.Engine, svcCtx)
	s.initApiSixConfig(svcCtx.Db)

	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		log.Info("Close api Server")
		s.genericApiServer.Close()
		return nil
	}))

	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}
	return s.genericApiServer.Run()
}

// buildGenericConfig 根据命令行参数创建 GenericAPIServer 配置
func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, lastErr error) {
	genericConfig = genericapiserver.NewConfig()
	if lastErr = cfg.ServerOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	return
}

// initDbStore 数据库连接初始化
func (s *apiServer) initDbStore() {
	var err error
	switch s.dbOption.Type { // 根据连接数据库的类型，进行不同驱动的适配
	case "postgresql":
		_, err = dal.GetPostgresFactoryOr(s.dbOption)
		if err != nil {
			log.Panicf("init postgres store failed: %s", err.Error())
			return
		}
		s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
			log.Infof("Close Postgres db")
			dbIns, _ := dal.GetPostgresFactoryOr(nil)
			return dbIns.CloseDb()
		}))
	case "mysql":
		_, err := dal.GetMySQLFactoryOr(s.dbOption)
		if err != nil {
			log.Errorf("init mysql store failed: %s", err.Error())
			return
		}
		s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
			log.Infof("Close MySQL db")
			dbIns, _ := dal.GetMySQLFactoryOr(nil)
			return dbIns.CloseDb()
		}))
	}
	return
}

// initApiSixConfig 项目启动首先初始化 apisix config
func (s *apiServer) initApiSixConfig(db *gorm.DB) {
	err := apisix.GenConfig(db)
	if err != nil {
		log.Panicf("init apisix config failed: %s", err.Error())
		return
	}
}
