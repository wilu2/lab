package options

import (
	"github.com/spf13/pflag"
	"gitlab.intsig.net/textin-gateway/internal/pkg/server"
)

type ServerOptions struct {
	Version     string   `json:"version" mapstructure:"version"`
	Host        string   `json:"host" mapstructure:"host"`
	Port        int      `json:"port" mapstructure:"port"`
	Mode        string   `json:"mode" mapstructure:"mode"`
	Domain      string   `json:"domain" mapstructure:"domain"`
	Health      bool     `json:"health"     mapstructure:"health"`
	Language    string   `json:"language" mapstructure:"language"`
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

func NewServerRunOptions() *ServerOptions {
	defaults := server.NewConfig()

	return &ServerOptions{
		Version:     "",
		Host:        defaults.Host,
		Port:        defaults.Port,
		Mode:        defaults.Mode,
		Domain:      "127.0.0.1:8080",
		Health:      defaults.Health,
		Language:    defaults.Language,
		Middlewares: defaults.Middlewares,
	}
}

// ApplyTo 从命令行读取到 Config 中
func (s *ServerOptions) ApplyTo(c *server.Config) error {
	c.Version = s.Version
	c.Host = s.Host
	c.Port = s.Port
	c.Mode = s.Mode
	c.Health = s.Health
	c.Middlewares = s.Middlewares
	c.Language = s.Language

	return nil
}

// Validate 添加参数校验
func (s *ServerOptions) Validate() []error {
	var errors []error
	return errors
}

// AddFlags 将特定APIServer的标志添加到指定的标志集中
func (s *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Host, "server.host", s.Host, "服务监听地址")
	fs.IntVar(&s.Port, "server.port", s.Port, "服务监听端口")
	fs.StringVar(&s.Mode, "server.mode", s.Mode, "服务的启动模式: debug, test, release.")
	fs.BoolVar(&s.Health, "server.health", s.Health, "添加服务监测路由 /health")
	fs.StringSliceVar(&s.Middlewares, "server.middlewares", s.Middlewares, "允许的中间件列表，逗号分隔，如果为空，则使用默认中间件")
	fs.StringSliceVar(&s.Middlewares, "server.language", s.Middlewares, "国际化语言: en, zh")
}
