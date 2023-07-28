package log

import (
	"fmt"
	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

const (
	flagLevel             = "log.level"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagFormat            = "log.format"
	flagEnableColor       = "log.enable-color"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagDevelopment       = "log.development"
	flagName              = "log.name"

	consoleFormat = "console"
	jsonFormat    = "json"
)

// Options 包含与日志相关的配置项
type Options struct {
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	Level             string   `json:"level"              mapstructure:"level"`
	Format            string   `json:"format"             mapstructure:"format"`
	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	EnableColor       bool     `json:"enable-color"       mapstructure:"enable-color"`
	Development       bool     `json:"development"        mapstructure:"development"`
	Name              string   `json:"name"               mapstructure:"name"`
}

// NewOptions 创建具有默认参数的 Options 对象
func NewOptions() *Options {
	return &Options{
		Level:             zapcore.InfoLevel.String(),
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            consoleFormat,
		EnableColor:       false,
		Development:       false,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}
}

// Validate 校验输入的配置项是否合法
func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

// AddFlags 将日志的标志添加到指定的FlagSet对象
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Name, flagName, o.Name, "Logger的名字")
	fs.StringVar(&o.Format, flagFormat, o.Format, "支持的日志输出格式, 目前支持console, json")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "是否开启颜色输出, true, false")
	fs.StringVar(&o.Level, flagLevel, o.Level, "最小日志级别, debug, info, warn, error, dpanic, panic, fatal")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "日志中显示调用日志所在的文件、函数和行号")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace, o.DisableStacktrace, "是否再panic及以上级别禁止打印堆栈信息")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "支持输出到多个输出，逗号分开, 支持输出到标准输出和文件")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "zap内部(非业务)错误日志输出路径, 支持多个输出")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"是否是开发模式。如果是开发模式，会对DPanicLevel进行堆栈跟踪",
	)
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}
