package options

import (
	"github.com/spf13/pflag"
	"gitlab.intsig.net/textin-gateway/internal/pkg/server"
)

// FeatureOptions API服务器功能相关的配置项
type FeatureOptions struct {
	EnableProfiling bool `json:"profiling" mapstructure:"profiling"`
	EnableMetrics   bool `json:"metrics" mapstructure:"metrics"`
}

func NewFeatureOptions() *FeatureOptions {
	defaults := server.NewConfig()

	return &FeatureOptions{
		EnableMetrics:   defaults.EnableMetrics,
		EnableProfiling: defaults.EnableProfiling,
	}
}

func (o *FeatureOptions) ApplyTo(c *server.Config) error {
	c.EnableProfiling = o.EnableProfiling
	c.EnableMetrics = o.EnableMetrics

	return nil
}

func (o *FeatureOptions) Validate() []error {
	return []error{}
}

func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.BoolVar(&o.EnableProfiling, "feature.profiling", o.EnableProfiling,
		"通过 web 界面启用分析 host:port/debug/pprof/")
	fs.BoolVar(&o.EnableMetrics, "feature.enable-metrics", o.EnableMetrics,
		"指标监测 /metrics")
}
