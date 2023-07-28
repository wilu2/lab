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

type postgresRepo struct {
	Db *gorm.DB
}

func (d *postgresRepo) GetDb() *gorm.DB {
	return d.Db
}

func (d *postgresRepo) CloseDb() error {
	db, err := d.Db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

var (
	postgresFactory Repo
	once            sync.Once
)

// GetMySQLFactoryOr create mysql factory with the given config.
func GetPostgresFactoryOr(opts *genericOptions.DbSQLOptions) (Repo, error) {
	if opts == nil && postgresFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}
	var err error
	var dbIns *gorm.DB
	mysqlOnce.Do(func() {
		options := &db.PostgresOptions{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			Port:                  opts.Port,
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

		dbIns, err = db.NewPostgres(options)

		// uncomment the following line if you need auto migration the given models
		// not suggested in production environment.
		// migrateDatabase(dbIns)
		postgresFactory = &postgresRepo{dbIns}
	})

	if err != nil {
		return postgresFactory, fmt.Errorf("failed to get postgres db")
	}

	return postgresFactory, nil
}
