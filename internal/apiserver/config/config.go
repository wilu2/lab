package config

import (
	"github.com/spf13/viper"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/options"
)

type Config struct {
	*options.Options
	Apisix *ApisixConfig `json:"apisix" mapstructure:"apisix"`
}

// CreateConfigFromOptions 创建一个配置实例
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	conf := &Config{
		Options: opts,
		Apisix:  nil,
	}
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// EmailConfig 邮件配置信息
type ApisixConfig struct {
	BaseUrl string `json:"base-url,omitempty" mapstructure:"base-url"`
	ApiKey  string `json:"api-key,omitempty" mapstructure:"api-key"`
}
