package app

import (
	cliflag "gitlab.intsig.net/textin-gateway/pkg/cli"
)

// CliOptions 抽象用于从命令行读取参数的配置选项
type CliOptions interface {
	Flags() (fss cliflag.NamedFlagSets)
	Validate() []error // 校验参数合法性
}

// ConfigurableOptions abstracts configuration options for reading parameters
// from a configuration file.
type ConfigurableOptions interface {
	// ApplyFlags parsing parameters from the command line or configuration file
	// to the options instance.
	ApplyFlags() []error
}

// CompleteAbleOptions 抽象出可以补全的选项
type CompleteAbleOptions interface {
	Complete() error
}

// PrintableOptions 抽象出可以打印的选项
type PrintableOptions interface {
	String() string
}
