package options

import (
	genericOptions "gitlab.intsig.net/textin-gateway/internal/pkg/options"
	cliflag "gitlab.intsig.net/textin-gateway/pkg/cli"
	"gitlab.intsig.net/textin-gateway/pkg/log"
)

// Options 输入命令配置参数
type Options struct {
	ServerOptions     *genericOptions.ServerOptions     `json:"server"   mapstructure:"server"`          // 基础 Server 配置
	ApisixOptions     *genericOptions.ApisixOptions     `json:"control"  mapstructure:"control"`         // Control-Panel 配置
	DbSQLOptions      *genericOptions.DbSQLOptions      `json:"db"    mapstructure:"db"`                 // Postgres 配置
	ClickHouseOptions *genericOptions.ClickHouseOptions `json:"clickhouse"    mapstructure:"clickhouse"` // ClickHouse 配置
	FeatureOptions    *genericOptions.FeatureOptions    `json:"feature"  mapstructure:"feature"`         // 性能、指标监测功能
	Log               *log.Options                      `json:"log"      mapstructure:"log"`             // 日志配置
}

func (o Options) Flags() (fss cliflag.NamedFlagSets) {
	o.ServerOptions.AddFlags(fss.FlagSet("server"))
	o.ApisixOptions.AddFlags(fss.FlagSet("apisix"))
	o.DbSQLOptions.AddFlags(fss.FlagSet("db"))
	o.ClickHouseOptions.AddFlags(fss.FlagSet("clickhouse"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.Log.AddFlags(fss.FlagSet("log"))
	return fss
}

func NewOptions() *Options {
	return &Options{
		ServerOptions:     genericOptions.NewServerRunOptions(),
		ApisixOptions:     genericOptions.NewApisixOptions(),
		DbSQLOptions:      genericOptions.NewDbSQLOptions(),
		ClickHouseOptions: genericOptions.NewClickHouseOptions(),
		FeatureOptions:    genericOptions.NewFeatureOptions(),
		Log:               log.NewOptions(),
	}
}

// Validate 校验每个命令分组的输入是否正确
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.ServerOptions.Validate()...)
	errs = append(errs, o.ApisixOptions.Validate()...)
	errs = append(errs, o.DbSQLOptions.Validate()...)
	errs = append(errs, o.ClickHouseOptions.Validate()...)
	errs = append(errs, o.FeatureOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return errs
}
