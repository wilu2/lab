package app

import (
	"fmt"

	"gitlab.intsig.net/textin-gateway/pkg/log"

	"github.com/fatih/color"

	"os"

	cliflag "gitlab.intsig.net/textin-gateway/pkg/cli"

	"github.com/marmotedu/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	progressMessage = color.GreenString("==>")
)

// App 应用程序的主要结构
type App struct {
	name        string     // APP 简短描述
	basename    string     // 二进制文件名称
	description string     // 应用详细描述
	options     CliOptions // 应用命令行选项
	runFunc     RunFunc    // 启动回调函数
	noConfig    bool       // 是否需要配置文件
	silence     bool       // 是否需要打印启动日志
	commands    []*Command
	args        cobra.PositionalArgs // 校验输入命令
	cmd         *cobra.Command
}

// Run 启动应用程序
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// Option 用于初始化应用程序结构的可选参数
type Option func(*App)

// WithOptions 从命令行读取或从配置文件读取参数
func WithOptions(opt CliOptions) Option {
	return func(app *App) {
		app.options = opt
	}
}

// WithDescription App 描述
func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

// RunFunc 定义应用程序的启动回调函数
type RunFunc func(basename string) error

// WithRunFunc 设置应用程序启动回调函数选项
func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

// NewApp 动态配置 APP
func NewApp(name string, basename string, opts ...Option) *App {
	app := &App{
		name:     name,
		basename: basename,
	}
	for _, o := range opts {
		o(app)
	}
	app.buildCommand()
	return app
}

// buildCommand 将启动命令参数映射到 cobra
func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:           a.basename,
		Short:         a.name,
		Long:          a.description,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	cliflag.InitFlags(cmd.Flags())

	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
	}

	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	var namedFlagSets cliflag.NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	globalFlagSet := namedFlagSets.FlagSet("global") // 设置全局命令
	cliflag.AddHelpFlags(globalFlagSet, cmd.Name())
	if !a.noConfig {
		addConfigFlag(globalFlagSet, a.basename)
	}

	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	a.cmd = &cmd
}

// runCommand 启动程序入口
func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	cliflag.PrintFlags(cmd.Flags())

	if !a.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return nil
		}
		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	if !a.silence {
		log.Infof("%v Starting %s ...", progressMessage, a.name)
		if !a.noConfig {
			log.Infof("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
		}
	}

	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}

	if a.runFunc != nil {
		return a.runFunc(a.basename) // 启动程序
	}

	return nil
}

// 调用字段检查规则，或增加打印方法
func (a *App) applyOptionRules() error {
	if errs := a.options.Validate(); len(errs) != 0 {
		return errors.NewAggregate(errs)
	}
	return nil
}

// WithDefaultValidArgs 输入参数校验
func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

// printWorkingDir 打印工作目录
func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
}
