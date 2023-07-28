# TextIn Gateway 项目

部署文档：https://gitlab.intsig.net/textin-gateway/textin-gateway-manifest

## 项目介绍

基于 `apisix` 做的二次功能开发，以及定制一些插件功能。
`apisix` 使用了 `Stand-alone` 模式，去除掉了 `etcd` 的依赖。
主要实现了 `apisix` 两个插件功能：
1. 将请求体和返回值保存在本地。
2. 将请求的历史记录保存在数据库的 `access_log` 表中，在访问统计页面查看。

### 功能分支

* `develop`：主要功能开发分支，以及对外提供的版本。
* `feature/bbw`：功能定制分支，主要有用户组，`sso` 登陆功能
* `feature/jianfa`: 访问统计日志，添加 `serviceName`

## 项目环境

基于 `golang` 的 `1.18` 版本

### 项目目录介绍
```shell
├── api
│  ├── ocr                   # ocr 结果的样例
│  └── swagger               # openapi 文件
├── build
│  └── docker                # 项目的 dockerfile 文件
├── cmd
│  ├── api-server            # 项目启动入口
│  └── gentool               # gentool 生成 GORM 的功能
├── configs                  # 配置文件
│  ├── api-server.yaml       # 项目配置文件信息
│  ├── init-route-base.yaml  # apisix 配置文件
│  └── vars.go
├── database
│  └── postgres              # postgresql 表创建以及更新语句
├── Dockerfile               # 项目的 Dockerfile 文件
├── go.mod                   # golang mod
├── go.sum
├── internal
│  ├── apiserver
│  │  ├── app.go
│  │  ├── code              # 业务代码 code 定义
│  │  ├── config            # 解析config 配置文件
│  │  ├── consts            # consts 常量定义
│  │  ├── dal               # dal 数据库的相关操作
│  │  │  ├── model          # 数据模型或者结构体
│  │  │  └── query          # gorm 定义的 query 语句
│  │  ├── handler           # handler 层主要是获取参数与参数校验的操作
│  │  ├── logic             # logic 层是完成功能逻辑
│  │  ├── middleware        # gin 中间件定义
│  │  ├── options
│  │  ├── response          # response 处理逻辑
│  │  ├── routes            # 定义 api router 入口
│  │  ├── server.go
│  │  ├── svc               # svc 资源定义
│  │  └── types             # 定义 Request Response 请求结构体
│  └── pkg
│     ├── factory           # mock apisix 信息
│     ├── jwtgen            # jwt token 生成逻辑
│     ├── middleware        # 一些常用中间件
│     ├── options           # viper 读取环境变量以及配置文件
│     ├── optlogs           # 操作日志逻辑
│     ├── server
│     ├── session
│     └── utils             # 工具库函数
├── pkg                     # 三方 pkg
│  ├── apisix               # apisix 的操作逻辑
│  ├── app
│  ├── db                   # 数据库连接逻辑
│  ├── encrypt              # 密码工具
│  ├── errors               # error 处理逻辑
│  ├── log                  # 日志模块
│  ├── md5                  # md5 处理逻辑
│  ├── openapi
│  ├── shutdown             # 项目 shotdown 清理逻辑
│  ├── sort                 # 排序规则实现
│  ├── stringx              # 字符串操作逻辑
│  └── validation           # validate 校验数据
└── README.md
```

### 用户权限

* 超管: 拥有所有数据权限和页面权限
* 管理员: 没有操作日志页面权限
* 普通用户: 只能查询，没有操作权限，可以访问体验中心和访问统计，平台概览和渠道列表

## 项目工具
### 工具安装

1. `go` 添加私库 `go env -w GOPRIVATE=gitlab.intsig.net`
2. `ginctl` 下载 `go install gitlab.intsig.net/wenbo_hu/ginctl`
3. `gorm-gen` 下载 `go install gorm.io/gen/tools/gentool@latest`

* [ginctl](https://gitlab.intsig.net/wenbo_hu/ginctl)：根据定义 `api` 文件生成 `swagger` 文档，和路由(routes)，请求参数(types)，资源处理(handler)，业务代码(logic)。
* [gen](https://github.com/go-gorm/gen/blob/master/README.ZH_CN.md)：指定数据库地址和表，生成 `model` 和 `query` 语句。
* [pre-commit](https://pre-commit.com/) ：用于代码提交规范代码检查

### 工具使用

1. 根据 `api` 生成 `swagger` 文档：`ginctl swagger -a ./internal/apiserver/api/user.api -d ./api -f user_swagger.json`
2. 根据 `api` 生成服务代码：`ginctl -a ./internal/apiserver/api/textin.api -d ./internal/apiserver`
3. 生成数据库： `model`：`go run ./gentool.go -db mysql -dsn "postgres:postgres@tcp(192.168.60.118:3306)/gateway" -outPath ../../internal/apiserver/dal/query`
* 项目启动；`go run ./cmd/api-server/apiserver.go -c ./configs/api-server.yaml`
