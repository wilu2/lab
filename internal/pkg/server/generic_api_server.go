package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"gitlab.intsig.net/textin-gateway/internal/pkg/middleware"
	"gitlab.intsig.net/textin-gateway/pkg/log"
	"golang.org/x/sync/errgroup"
)

// GenericApiServer gin.Engine Api包装器
type GenericApiServer struct {
	*gin.Engine
	Address         string        // 监听地址端口
	health          bool          // 健康检查
	middlewares     []string      // 配置文件开启中间件
	ShutdownTimeout time.Duration // 服务器关闭超时
	secureServer    *http.Server
	enableMetrics   bool
	enableProfiling bool // 是否开启 pprof 处理程序
}

// initGenericApiServer 完成 ApiServer 初始化配置
func initGenericApiServer(s *GenericApiServer) {
	s.Setup()              // 做一些 gin.Engine 初始化工作
	s.InstallMiddlewares() // 添加默认中间件
	s.InstallAPIs()        // 添加一些 api
}

// Setup gin.engine 初始化工作
func (s *GenericApiServer) Setup() {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

// InstallMiddlewares 根据配置文件或传参数 安装全局 middlewares
func (s *GenericApiServer) InstallMiddlewares() {
	for _, m := range s.middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)
			continue
		}

		log.Infof("install middleware: %s", m)
		s.Use(mw)
	}
}

// InstallAPIs 根据配置或传参开启一些 api
func (s *GenericApiServer) InstallAPIs() {
	if s.health {
		s.GET("/health", func(c *gin.Context) {
			core.WriteResponse(c, nil, map[string]string{"status": "ok"})
		})
	}

	if s.enableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.Engine)
	}

	if s.enableProfiling {
		pprof.Register(s.Engine)
	}
}

// Run 开启服务监听端口
func (s *GenericApiServer) Run() error {
	s.secureServer = &http.Server{
		Addr:    s.Address,
		Handler: s,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	} // 自定义HTTP配置模式

	var eg errgroup.Group

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	eg.Go(func() error {
		log.Infof("Start to listening the incoming requests on http address: %s", s.Address)

		if err := s.secureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.Address)

		return nil
	})
	// Ping the server to make sure the router is working.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.health {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil

}

// Close 优雅地关闭api服务器
func (s *GenericApiServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.secureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown secure server failed: %s", err.Error())
	}
}

func (s *GenericApiServer) ping(ctx context.Context) error {
	return nil
}
