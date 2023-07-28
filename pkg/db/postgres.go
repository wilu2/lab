package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Options defines optsions for postgres database.
type PostgresOptions struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	Port                  int
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
	Logger                logger.Interface
}

// New create a new gorm db instance with the given options.
func NewPostgres(opts *PostgresOptions) (*gorm.DB, error) {
	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC`,
		opts.Host,
		opts.Username,
		opts.Password,
		opts.Database,
		opts.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
