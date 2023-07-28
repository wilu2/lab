package apisix

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/spf13/viper"
	"gitlab.intsig.net/textin-gateway/configs"
	"gitlab.intsig.net/textin-gateway/pkg/log"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

var (
	once             sync.Once
	routeInstance    []string
	upstreamInstance []string
	genConfigLock    sync.Mutex // 写入 yaml 文件时的互斥锁
)

type Config struct {
	Routes    []Route    `yaml:"routes"`
	Upstreams []Upstream `yaml:"upstreams"`
}

const tmpl = `routes:
{{- range .Routes }}
  -
{{ . }}
{{- end }}

upstreams:
{{- range .Streams }}
  -
{{ . }}
{{- end }}

global_rules:
  -
    id: 1
    plugins:
      request-id:
        header_name: X-Request-Id
        algorithm: uuid
        disable: false
        include_in_response: true

#END
`

// GenConfig 根据数据库中路由规则生成 apisix 配置文件
func GenConfig(db *gorm.DB) error {
	var (
		buf        bytes.Buffer
		configPath = viper.GetString("control.config")
	)
	routes, streams, err := GetConfigInfo(db)
	data := struct {
		Routes  []string
		Streams []string
	}{
		Routes:  routes,
		Streams: streams,
	}
	tpl := template.Must(template.New("tmpl").Parse(tmpl))
	if err := tpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("执行模板失败 %v", err)
	}
	genConfigLock.Lock()
	defer func() {
		genConfigLock.Unlock()
		if r := recover(); r != nil {
			log.Errorf("发生 panic：%v\n", r)
		}
	}()
	if err := ioutil.WriteFile(configPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入文件失败 %v", err)
	}
	return err
}

// GetConfigInfo 通过环境变量和数据库获得 Route 和 Upstream
func GetConfigInfo(db *gorm.DB) ([]string, []string, error) {
	var (
		routes        []string
		streams       []string
		routesChan    = make(chan []string)
		upstreamsChan = make(chan []string)
	)
	once.Do(func() {
		routeInstance, upstreamInstance = InitConfig()
	})

	// 同时获取路由和上游信息
	go func() {
		routes, err := AllRouteYaml(db)
		if err != nil {
			log.Errorf("获取数据库 route 配置信息失败 %v", err.Error())
		}
		routesChan <- routes
	}()
	go func() {
		streams, err := AllStreamYaml(db)
		if err != nil {
			log.Errorf("获取数据库 upstream 配置信息失败 %v", err.Error())
		}
		upstreamsChan <- streams
	}()

	select {
	case routes = <-routesChan:
		routes = append(routes, routeInstance...)
	case <-time.After(time.Second):
		log.Error("获取路由信息超时")
		return nil, nil, errors.New("获取路由信息超时")
	}

	select {
	case streams = <-upstreamsChan:
		streams = append(streams, upstreamInstance...)
	case <-time.After(time.Second):
		log.Error("获取上游信息超时")
		return nil, nil, errors.New("获取上游信息超时")
	}
	return routes, streams, nil
}

// InitConfig 通过读取环境变量，将 训练平台，管理平台，合同生成路由注册信息
func InitConfig() (routeInfoList []string, upstreamInfoList []string) {
	routeAddressVars := make(map[string]string)
	envVars := os.Environ()
	// 遍历所有环境变量，检查是否以 ROUTE 开头并以 ADDRESS 结尾
	for _, envVar := range envVars {
		envVarList := strings.Split(envVar, "=")
		if strings.HasPrefix(envVarList[0], "ROUTE") && strings.HasSuffix(envVarList[0], "ADDRESS") {
			routeAddressVars[envVarList[0]] = envVarList[1]
		}
	}
	var config Config
	var configByte []byte

	file, err := os.Open("/app/configs/init-route.yaml")
	if err != nil {
		log.Info("使用默认 apisix route 文件")
	}
	defer file.Close()
	if err == nil {
		configByte, err = ioutil.ReadAll(file)
		if len(configByte) == 0 {
			log.Info("初始化 route 文件为空，使用默认 apisix route 文件")
			err = errors.New("文件为空")
		}
	}

	if err != nil {
		_ = yaml.Unmarshal(configs.InitRouteConfig, &config)
	} else {
		_ = yaml.Unmarshal(configByte, &config)
	}

	flagList := []int{}
	for i, stream := range config.Upstreams {
		name := "ROUTE_" + stream.Name + "_ADDRESS"
		if address, ok := routeAddressVars[name]; ok { // 如果配置了环境变量则使用，配置的地址
			log.Infof("%s 路由自定义地址 %s", stream.Name, address)
			config.Upstreams[i].Nodes = map[string]int{address: 1}
		}
		flagList = append(flagList, i)
	}

	for _, i := range flagList {
		routeStr, _ := yaml.Marshal(config.Routes[i])
		repRouteStr := "    " + strings.Replace(string(routeStr), "\n", "\n    ", -1)
		routeInfoList = append(routeInfoList, repRouteStr)
		streamStr, _ := yaml.Marshal(config.Upstreams[i])
		repStreamStr := "    " + strings.Replace(string(streamStr), "\n", "\n    ", -1)
		upstreamInfoList = append(upstreamInfoList, repStreamStr)
	}
	return routeInfoList, upstreamInfoList
}
