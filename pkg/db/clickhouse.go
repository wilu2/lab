package db

import (
	"fmt"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Options defines options for clickhouse database.
type ClickHouseOptions struct {
	Host                  string
	Port                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
	Logger                logger.Interface
}

// New create a new gorm db instance with the given options.
func NewClickHouse(opts *ClickHouseOptions) (*gorm.DB, error) {
	dsn := fmt.Sprintf(`clickhouse://%s:%s@%s:%s/%s?dial_timeout=10s&read_timeout=20s`,
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Port,
		opts.Database,
	)

	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{
		Logger: opts.Logger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}
