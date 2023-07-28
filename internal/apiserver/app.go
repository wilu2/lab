package apiserver

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/config"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/options"
	"gitlab.intsig.net/textin-gateway/pkg/app"
	"gitlab.intsig.net/textin-gateway/pkg/log"
	"gitlab.intsig.net/textin-gateway/pkg/validation"
)

const desc = "TextIn Gateway Backend"

func NewApp(basename string) *app.App {
	opts := options.NewOptions() // 初始化配置参数
	application := app.NewApp(
		"Api Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(desc),
		app.WithDefaultValidArgs(), // 命令行非选项参数的默认校验逻辑
		app.WithRunFunc(run(opts)), // 应用启动函数
	)
	return application
}

// run 函数启动入口
func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log) // 初始化日志参数
		defer log.Flush()
		err := validation.Init(opts.ServerOptions.Language)
		if err != nil {
			return err
		}
		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}
		server, err := createAPIServer(cfg)
		if err != nil {
			return err
		}
		return server.Run()
	}
}
