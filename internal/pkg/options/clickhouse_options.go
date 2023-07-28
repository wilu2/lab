package options

import (
	"time"

	"github.com/spf13/pflag"
	"gitlab.intsig.net/textin-gateway/internal/pkg/server"
	"gitlab.intsig.net/textin-gateway/pkg/db"
	"gorm.io/gorm"
)

type ClickHouseOptions struct {
	Host                  string        `json:"host,omitempty"                     mapstructure:"host"`
	Port                  string        `json:"port,omitempty"                     mapstructure:"port"`
	Username              string        `json:"username,omitempty"                 mapstructure:"username"`
	Password              string        `json:"password"                           mapstructure:"password"`
	Database              string        `json:"database"                           mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"     mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	LogLevel              int           `json:"log-level"                          mapstructure:"log-level"`
}

func NewClickHouseOptions() *ClickHouseOptions {
	return &ClickHouseOptions{
		Host:     "127.0.0.1",
		Port:     "9000",
		Username: "",
		Password: "",
		Database: "",
		LogLevel: 1, // Silent
	}
}

// Validate 添加参数校验
func (o *ClickHouseOptions) Validate() []error {
	var errors []error
	return errors
}

// AddFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet.
func (o *ClickHouseOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "clickhouse.host", o.Host, "ClickHouse 服务主机地址")
	fs.StringVar(&o.Host, "clickhouse.port", o.Port, "ClickHouse 服务主机端口")
	fs.StringVar(&o.Username, "clickhouse.username", o.Username, "ClickHouse 服务用户名")
	fs.StringVar(&o.Password, "clickhouse.password", o.Password, "ClickHouse 服务密码")
	fs.StringVar(&o.Database, "clickhouse.database", o.Database, "ClickHouse 数据库名称")
	fs.IntVar(&o.LogLevel, "clickhouse.log-mode", o.LogLevel, "指定 gorm 日志级别")
}

func (o *ClickHouseOptions) ApplyTo(c *server.Config) error {
	return nil
}

// NewClient 使用给定的配置创建 ClickHouse 实例
func (o *ClickHouseOptions) NewClient() (*gorm.DB, error) {
	opts := &db.ClickHouseOptions{
		Host:     o.Host,
		Port:     o.Port,
		Username: o.Username,
		Password: o.Password,
		Database: o.Database,
		LogLevel: o.LogLevel,
	}

	return db.NewClickHouse(opts)
}
