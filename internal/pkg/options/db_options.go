package options

import (
	"strconv"
	"time"

	"github.com/spf13/pflag"
	"gitlab.intsig.net/textin-gateway/internal/pkg/server"
)

type DbSQLOptions struct {
	Type                  string        `json:"type,omitempty" mapstructure:"type"`
	Host                  string        `json:"host,omitempty" mapstructure:"host"`
	Username              string        `json:"username,omitempty" mapstructure:"username"`
	Password              string        `json:"-" mapstructure:"password"`
	Database              string        `json:"database,omitempty" mapstructure:"database"`
	Port                  int           `json:"port,omitempty" mapstructure:"port"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty" mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty" mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	LogLevel              int           `json:"log-level" mapstructure:"log-level"`
}

func NewDbSQLOptions() *DbSQLOptions {
	return &DbSQLOptions{
		Type:                  "postgres",
		Host:                  "127.0.0.1",
		Username:              "",
		Password:              "",
		Database:              "",
		Port:                  5432,
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              1, // Silent
	}
}

// Validate 添加参数校验
func (o *DbSQLOptions) Validate() []error {
	var errors []error
	return errors
}

// AddFlags adds flags related to postgres storage for a specific APIServer to the specified FlagSet.
func (o *DbSQLOptions) AddFlags(fs *pflag.FlagSet) {
	port := strconv.Itoa(o.Port)
	fs.StringVar(&o.Host, "db.host", o.Host, "数据库服务主机地址")
	fs.StringVar(&o.Username, "db.username", o.Username, "数据库服务用户名")
	fs.StringVar(&o.Password, "db.password", o.Password, "数据库服务密码")
	fs.StringVar(&o.Database, "db.database", o.Database, "数据库名称")
	fs.StringVar(&port, "db.port", strconv.Itoa(o.Port), "数据库端口")
	fs.IntVar(&o.MaxIdleConnections, "db.max-idle-connections", o.MaxOpenConnections, ""+
		"允许连接到数据库的最大空闲连接数")
	fs.IntVar(&o.MaxOpenConnections, "db.max-open-connections", o.MaxOpenConnections, ""+
		"允许连接到数据库的最大打开连接数。")
	fs.DurationVar(&o.MaxConnectionLifeTime, "db.max-connection-life-time", o.MaxConnectionLifeTime, ""+
		"允许连接到数据库的最长连接生存时间")
	fs.IntVar(&o.LogLevel, "db.log-mode", o.LogLevel, "指定 gorm 日志级别")
}

func (o *DbSQLOptions) ApplyTo(c *server.Config) error {
	return nil
}
