package options

import (
	"github.com/spf13/pflag"
	"gitlab.intsig.net/textin-gateway/internal/pkg/server"
)

type ApisixOptions struct {
	Host          string `json:"host,omitempty" mapstructure:"host"`
	Port          string `json:"port" mapstructure:"port"`
	Apikey        string `json:"api-key,omitempty" mapstructure:"api-key"`
	BodyLoggerDir string `json:"body-logger-dir,omitempty" mapstructure:"body-logger-dir"`
	Config        string `json:"config,omitempty" mapstructure:"config"`
}

func NewApisixOptions() *ApisixOptions {
	return &ApisixOptions{
		Host:          "127.0.0.1",
		Port:          "9180",
		Apikey:        "edd1c9f034335f136f87ad84b625c8f1",
		BodyLoggerDir: "/tmp",
		Config:        "/tmp/apisix.yaml",
	}
}

// Validate 添加参数校验
func (o *ApisixOptions) Validate() []error {
	var errors []error
	return errors
}

// AddFlags adds flags related to postgres storage for a specific APIServer to the specified FlagSet.
func (o *ApisixOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "control.ip", o.Host, "Apisix 服务主机地址")
	fs.StringVar(&o.Port, "control.port", o.Port, "Apisix 服务端口")
	fs.StringVar(&o.Apikey, "control.api-key", o.Apikey, "Apisix Admin Key")
	fs.StringVar(&o.BodyLoggerDir, "control.body-logger-dir", o.BodyLoggerDir, "Apisix Body Logger Dir")
	fs.StringVar(&o.Config, "control.config", o.Config, "Apisix Config")
}

func (o *ApisixOptions) ApplyTo(c *server.Config) error {
	return nil
}
