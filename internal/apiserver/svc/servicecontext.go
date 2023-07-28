package svc

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/config"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/dal"
	"gorm.io/gorm"
)

var (
	svcInstance *ServiceContext
)

type ServiceContext struct {
	Config *config.Config
	Db     *gorm.DB
}

func NewServiceContext(cfg *config.Config) *ServiceContext {
	var dbIns dal.Repo
	switch cfg.DbSQLOptions.Type { // 根据连接数据库的类型，进行不同驱动的适配
	case "postgresql":
		dbIns, _ = dal.GetPostgresFactoryOr(nil)
	case "mysql":
		dbIns, _ = dal.GetMySQLFactoryOr(nil)
	}
	svcInstance = &ServiceContext{
		Config: cfg,
		Db:     dbIns.GetDb(),
	}
	return svcInstance
}

func GetSvc() *ServiceContext {
	return svcInstance
}
