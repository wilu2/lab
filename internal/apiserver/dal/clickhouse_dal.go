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

type ClickHouseRepo interface {
	GetClickHouse() *gorm.DB
	CloseClickHouse() error
}

type clickHouseRepo struct {
	Db *gorm.DB
}

func (d *clickHouseRepo) GetClickHouse() *gorm.DB {
	return d.Db
}

func (d *clickHouseRepo) CloseClickHouse() error {
	db, err := d.Db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

var (
	clickHouseFactory ClickHouseRepo
	clickHouseOnce    sync.Once
)

// GetMySQLFactoryOr create mysql factory with the given config.
func GetClickHouseFactoryOr(opts *genericOptions.ClickHouseOptions) (ClickHouseRepo, error) {
	if opts == nil && clickHouseFactory == nil {
		return nil, fmt.Errorf("failed to get clickhouse store factory")
	}
	// fmt.Sprintln(opts)
	var err error
	var dbIns *gorm.DB
	clickHouseOnce.Do(func() {
		options := &db.ClickHouseOptions{
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
		dbIns, err = db.NewClickHouse(options)

		// uncomment the following line if you need auto migration the given models
		// not suggested in production environment.
		// migrateDatabase(dbIns)
		clickHouseFactory = &clickHouseRepo{dbIns}
	})

	if err != nil {
		return clickHouseFactory, fmt.Errorf("failed to get clickhouse db")
	}

	return clickHouseFactory, nil
}
