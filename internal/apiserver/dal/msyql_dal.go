package dal

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/gorm/logger"

	"github.com/pkg/errors"
	genericOptions "gitlab.intsig.net/textin-gateway/internal/pkg/options"
	"gitlab.intsig.net/textin-gateway/pkg/db"
	"gorm.io/gorm"
)

type Repo interface {
	GetDb() *gorm.DB
	CloseDb() error
}

type dbRepo struct {
	Db *gorm.DB
}

func (d *dbRepo) GetDb() *gorm.DB {
	return d.Db
}

func (d *dbRepo) CloseDb() error {
	db, err := d.Db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

var (
	mysqlFactory Repo
	mysqlOnce    sync.Once
)

// GetMySQLFactoryOr create mysql factory with the given config.
func GetMySQLFactoryOr(opts *genericOptions.DbSQLOptions) (Repo, error) {
	if opts == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	var dbIns *gorm.DB
	mysqlOnce.Do(func() {
		options := &db.Options{
			Host:                  opts.Host,
			Port:                  opts.Port,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.LogLevel(opts.LogLevel),
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			}),
		}
		dbIns, err = db.New(options)

		// uncomment the following line if you need auto migration the given models
		// not suggested in production environment.
		// migrateDatabase(dbIns)
		mysqlFactory = &dbRepo{dbIns}
	})

	if err != nil {
		return mysqlFactory, fmt.Errorf("failed to get mysql db")
	}

	return mysqlFactory, nil
}
