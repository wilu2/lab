package apisix

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type BaseInfo struct {
	ID         interface{} `json:"id" yaml:"id"`
	CreateTime int64       `json:"create_time,omitempty" yaml:"create_time,omitempty"`
	UpdateTime int64       `json:"update_time,omitempty" yaml:"update_time,omitempty"`
}

func (info *BaseInfo) GetBaseInfo() *BaseInfo {
	return info
}

func (info *BaseInfo) Creating() {
	if info.ID == nil {
		info.ID = GetFlakeUidStr()
	} else {
		// convert to string if it's not
		if reflect.TypeOf(info.ID).String() != "string" {
			info.ID = InterfaceToString(info.ID)
		}
	}
	info.CreateTime = time.Now().Unix()
	info.UpdateTime = time.Now().Unix()
}

func (info *BaseInfo) Updating(storedInfo *BaseInfo) {
	info.ID = storedInfo.ID
	info.CreateTime = storedInfo.CreateTime
	info.UpdateTime = time.Now().Unix()
}

func (info *BaseInfo) KeyCompat(key string) {
	if info.ID == nil && key != "" {
		info.ID = key
	}
}

type Status uint8

// swagger:model Route
type Route struct {
	BaseInfo        `yaml:",inline"`
	URI             string                 `json:"uri,omitempty" yaml:"uri,omitempty"`
	Uris            []string               `json:"uris,omitempty" yaml:"uris,omitempty"`
	Name            string                 `json:"name" validate:"max=50" yaml:"name"`
	Desc            string                 `json:"desc,omitempty" validate:"max=256" yaml:"desc,omitempty"`
	Priority        int                    `json:"priority,omitempty" yaml:"priority,omitempty"`
	EnableWebSocket bool                   `json:"enable_websocket,omitempty" yaml:"enable_websocket,omitempty"`
	Methods         []string               `json:"methods,omitempty" yaml:"methods,omitempty"`
	Host            string                 `json:"host,omitempty" yaml:"host,omitempty"`
	Hosts           []string               `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	RemoteAddr      string                 `json:"remote_addr,omitempty" yaml:"remote_addr,omitempty"`
	RemoteAddrs     []string               `json:"remote_addrs,omitempty" yaml:"remote_addrs,omitempty"`
	Vars            []interface{}          `json:"vars,omitempty" yaml:"vars,omitempty"`
	FilterFunc      string                 `json:"filter_func,omitempty" yaml:"filter_func,omitempty"`
	Script          interface{}            `json:"script,omitempty" yaml:"script,omitempty"`
	ScriptID        interface{}            `json:"script_id,omitempty" yaml:"script_id,omitempty"` // For debug and optimization(cache), currently same as Route's ID
	Plugins         map[string]interface{} `json:"plugins" yaml:"plugins,omitempty"`
	Upstream        *UpstreamDef           `json:"upstream,omitempty" yaml:"upstream,omitempty"`
	ServiceID       interface{}            `json:"service_id,omitempty" yaml:"service_id,omitempty"`
	UpstreamID      interface{}            `json:"upstream_id,omitempty" yaml:"upstream_id,omitempty"`
	Labels          map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Status          Status                 `json:"status" yaml:"status"`
	ServiceProtocol string                 `json:"service_protocol,omitempty" yaml:"service_protocol,omitempty"`
	PluginConfigID  interface{}            `json:"plugin_config_id,omitempty" yaml:"plugin_config_id,omitempty"`
}

// --- structures for upstream start  ---
type TimeoutValue float32
type Timeout struct {
	Connect TimeoutValue `json:"connect,omitempty"`
	Send    TimeoutValue `json:"send,omitempty"`
	Read    TimeoutValue `json:"read,omitempty"`
}

type Node struct {
	Host     string      `json:"host,omitempty"`
	Port     int         `json:"port,omitempty"`
	Weight   int         `json:"weight"`
	Metadata interface{} `json:"metadata,omitempty"`
	Priority int         `json:"priority,omitempty"`
}

type K8sInfo struct {
	Namespace   string `json:"namespace,omitempty"`
	DeployName  string `json:"deploy_name,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	Port        int    `json:"port,omitempty"`
	BackendType string `json:"backend_type,omitempty"`
}

type Healthy struct {
	Interval     int   `json:"interval,omitempty"`
	HttpStatuses []int `json:"http_statuses,omitempty"`
	Successes    int   `json:"successes,omitempty"`
}

type UnHealthy struct {
	Interval     int   `json:"interval,omitempty"`
	HTTPStatuses []int `json:"http_statuses,omitempty"`
	TCPFailures  int   `json:"tcp_failures,omitempty"`
	Timeouts     int   `json:"timeouts,omitempty"`
	HTTPFailures int   `json:"http_failures,omitempty"`
}

type Active struct {
	Type                   string       `json:"type,omitempty"`
	Timeout                TimeoutValue `json:"timeout,omitempty"`
	Concurrency            int          `json:"concurrency,omitempty"`
	Host                   string       `json:"host,omitempty"`
	Port                   int          `json:"port,omitempty"`
	HTTPPath               string       `json:"http_path,omitempty"`
	HTTPSVerifyCertificate bool         `json:"https_verify_certificate,omitempty"`
	Healthy                Healthy      `json:"healthy,omitempty"`
	UnHealthy              UnHealthy    `json:"unhealthy,omitempty"`
	ReqHeaders             []string     `json:"req_headers,omitempty"`
}

type Passive struct {
	Type      string    `json:"type,omitempty"`
	Healthy   Healthy   `json:"healthy,omitempty"`
	UnHealthy UnHealthy `json:"unhealthy,omitempty"`
}

type HealthChecker struct {
	Active  Active  `json:"active,omitempty"`
	Passive Passive `json:"passive,omitempty"`
}

type UpstreamTLS struct {
	ClientCert string `json:"client_cert,omitempty"`
	ClientKey  string `json:"client_key,omitempty"`
}

type UpstreamKeepalivePool struct {
	IdleTimeout *TimeoutValue `json:"idle_timeout,omitempty"`
	Requests    int           `json:"requests,omitempty"`
	Size        int           `json:"size"`
}

type UpstreamDef struct {
	Nodes         interface{}            `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	Retries       *int                   `json:"retries,omitempty" yaml:"retries,omitempty"`
	Timeout       *Timeout               `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Type          string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Checks        interface{}            `json:"checks,omitempty" yaml:"checks,omitempty"`
	HashOn        string                 `json:"hash_on,omitempty" yaml:"hash_on,omitempty"`
	Key           string                 `json:"key,omitempty" yaml:"key,omitempty"`
	Scheme        string                 `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	DiscoveryType string                 `json:"discovery_type,omitempty" yaml:"discovery_type,omitempty"`
	DiscoveryArgs map[string]string      `json:"discovery_args,omitempty" yaml:"discovery_args,omitempty"`
	PassHost      string                 `json:"pass_host,omitempty" yaml:"pass_host,omitempty"`
	UpstreamHost  string                 `json:"upstream_host,omitempty" yaml:"upstream_host,omitempty"`
	Name          string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Desc          string                 `json:"desc,omitempty" yaml:"desc,omitempty"`
	ServiceName   string                 `json:"service_name,omitempty" yaml:"service_name,omitempty"`
	Labels        map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	TLS           *UpstreamTLS           `json:"tls,omitempty" yaml:"tls,omitempty"`
	KeepalivePool *UpstreamKeepalivePool `json:"keepalive_pool,omitempty" yaml:"keepalive_pool,omitempty"`
	RetryTimeout  TimeoutValue           `json:"retry_timeout,omitempty" yaml:"retry_timeout,omitempty"`
}

// swagger:model Upstream
type Upstream struct {
	BaseInfo    `yaml:",inline"`
	UpstreamDef `yaml:",inline"`
}

func (n Upstream) ConvertYAML() (string, error) {
	nodes, ok := n.Nodes.([]interface{})
	if !ok {
		// 如果不是，返回原始结构体
		context, err := yaml.Marshal(n)
		return string(context), err
	}
	nodeMap := make(map[string]int)
	for _, node := range nodes {
		// 对于每个node，将host与port组合成key，并将weight作为value存储到map中
		nodeInfo, ok := node.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid node format: %v", node)
		}
		port, ok := nodeInfo["port"].(float64)
		if !ok {
			return "", fmt.Errorf("invalid weight format: %v", nodeInfo["weight"])
		}
		weight, ok := nodeInfo["weight"].(float64)
		if !ok {
			return "", fmt.Errorf("invalid weight format: %v", nodeInfo["weight"])
		}
		key := fmt.Sprintf("%s:%d", nodeInfo["host"], int(port))
		nodeMap[key] = int(weight)
	}
	n.Nodes = nodeMap
	context, err := yaml.Marshal(n)
	contextStr := "    " + strings.Replace(string(context), "\n", "\n    ", -1)
	return contextStr, err
}

type UpstreamNameResponse struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
}

func (upstream *Upstream) Parse2NameResponse() (*UpstreamNameResponse, error) {
	nameResp := &UpstreamNameResponse{
		ID:   upstream.ID,
		Name: upstream.Name,
	}
	return nameResp, nil
}

// --- structures for upstream end  ---

// swagger:model Consumer
type Consumer struct {
	Username   string                 `json:"username"`
	Desc       string                 `json:"desc,omitempty"`
	Plugins    map[string]interface{} `json:"plugins,omitempty"`
	Labels     map[string]string      `json:"labels,omitempty"`
	CreateTime int64                  `json:"create_time,omitempty"`
	UpdateTime int64                  `json:"update_time,omitempty"`
}

type SSLClient struct {
	CA    string `json:"ca,omitempty"`
	Depth int    `json:"depth,omitempty"`
}

// swagger:model SSL
type SSL struct {
	BaseInfo
	Cert          string            `json:"cert,omitempty"`
	Key           string            `json:"key,omitempty"`
	Sni           string            `json:"sni,omitempty"`
	Snis          []string          `json:"snis,omitempty"`
	Certs         []string          `json:"certs,omitempty"`
	Keys          []string          `json:"keys,omitempty"`
	ExpTime       int64             `json:"exptime,omitempty"`
	Status        int               `json:"status"`
	ValidityStart int64             `json:"validity_start,omitempty"`
	ValidityEnd   int64             `json:"validity_end,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	Client        *SSLClient        `json:"client,omitempty"`
}

// swagger:model Service
type Service struct {
	BaseInfo
	Name            string                 `json:"name,omitempty"`
	Desc            string                 `json:"desc,omitempty"`
	Upstream        *UpstreamDef           `json:"upstream,omitempty"`
	UpstreamID      interface{}            `json:"upstream_id,omitempty"`
	Plugins         map[string]interface{} `json:"plugins,omitempty"`
	Script          string                 `json:"script,omitempty"`
	Labels          map[string]string      `json:"labels,omitempty"`
	EnableWebsocket bool                   `json:"enable_websocket,omitempty"`
	Hosts           []string               `json:"hosts,omitempty"`
}

type Script struct {
	ID     string      `json:"id"`
	Script interface{} `json:"script,omitempty"`
}

type RequestValidation struct {
	Type       string      `json:"type,omitempty"`
	Required   []string    `json:"required,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

// swagger:model GlobalPlugins
type GlobalPlugins struct {
	BaseInfo
	Plugins map[string]interface{} `json:"plugins"`
}

type ServerInfo struct {
	BaseInfo
	LastReportTime int64  `json:"last_report_time,omitempty"`
	UpTime         int64  `json:"up_time,omitempty"`
	BootTime       int64  `json:"boot_time,omitempty"`
	EtcdVersion    string `json:"etcd_version,omitempty"`
	Hostname       string `json:"hostname,omitempty"`
	Version        string `json:"version,omitempty"`
}

// swagger:model GlobalPlugins
type PluginConfig struct {
	BaseInfo
	Desc    string                 `json:"desc,omitempty" validate:"max=256"`
	Plugins map[string]interface{} `json:"plugins"`
	Labels  map[string]string      `json:"labels,omitempty"`
}

// swagger:model Proto
type Proto struct {
	BaseInfo
	Desc    string `json:"desc,omitempty"`
	Content string `json:"content"`
}

// swagger:model StreamRoute
type StreamRoute struct {
	BaseInfo
	Desc       string                 `json:"desc,omitempty"`
	RemoteAddr string                 `json:"remote_addr,omitempty"`
	ServerAddr string                 `json:"server_addr,omitempty"`
	ServerPort int                    `json:"server_port,omitempty"`
	SNI        string                 `json:"sni,omitempty"`
	Upstream   *UpstreamDef           `json:"upstream,omitempty"`
	UpstreamID interface{}            `json:"upstream_id,omitempty"`
	Plugins    map[string]interface{} `json:"plugins,omitempty"`
}

// swagger:model SystemConfig
type SystemConfig struct {
	ConfigName string                 `json:"config_name"`
	Desc       string                 `json:"desc,omitempty"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
	CreateTime int64                  `json:"create_time,omitempty"`
	UpdateTime int64                  `json:"update_time,omitempty"`
}
