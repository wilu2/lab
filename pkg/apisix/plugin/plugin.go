package plugin

type IPRestriction struct {
	WhiteList []string `json:"white_list"`
	BlackList []string `json:"black_list"`
	Message   string   `json:"message"`
}

type Cors struct {
	AllowCredential bool   `json:"allow_credential"  default:"false"`
	AllowHeaders    string `json:"allow_headers"  default:"*"`
	AllowMethods    string `json:"allow_methods"  default:"*"`
	AllowOrigins    string `json:"allow_origins"  default:"*"`
	Disable         bool   `json:"disable"  default:"true"`
	ExposeHeaders   string `json:"expose_headers"  default:"*"`
	MaxAge          int    `json:"max_age"  default:"5"`
}

type RequestId struct {
	Disable bool `json:"disable"  default:"true"`
}

type ProxyRewrite struct {
	Headers map[string]interface{} `json:"headers"`
	Host    string                 `json:"host"`
	Method  string                 `json:"method"`
	Scheme  string                 `json:"scheme"`
	URI     string                 `json:"uri"`
}

type ApiBreaker struct {
	BreakResponseBody    string                 `json:"break_response_body"`
	BreakResponseCode    int                    `json:"break_response_code"`
	BreakResponseHeaders []BreakResponseHeaders `json:"break_response_headers"`
	Disable              bool                   `json:"disable"`
	Healthy              HealthyStatus          `json:"healthy"`
	MaxBreakerSec        int                    `json:"max_breaker_sec"`
	Unhealthy            UnhealthyStatus        `json:"unhealthy"`
}

type LimitCount struct {
	Count                int      `json:"count"`
	TimeWindow           int      `json:"time_window"`
	KeyType              string   `json:"key_type"`
	Key                  string   `json:"key"`
	RejectedCode         int      `json:"rejected_code"`
	RejectedMsg          string   `json:"rejected_msg"`
	Policy               string   `json:"policy"`
	AllowDegradation     bool     `json:"allow_degradation"`
	ShowLimitQuotaHeader bool     `json:"show_limit_quota_header"`
	Group                string   `json:"group"`
	RedisHost            string   `json:"redis_host"`
	RedisPort            int      `json:"redis_port"`
	RedisPassword        string   `json:"redis_password"`
	RedisDataBase        int      `json:"redis_database"`
	RedisTimeout         int      `json:"redis_timeout"`
	RedisClusterNodes    []string `json:"redis_cluster_nodes"`
	RedisClusterName     string   `json:"redis_cluster_name"`
}

type ClientControl struct {
	MaxBodySize int `json:"max_body_size"` //Byte
}

type ClickHouseLogger struct {
	EndpointAddr  string   `json:"endpoint_addr"`
	EndpointAddrs []string `json:"endpoint_addrs"`
	Database      string   `json:"database"`
	Logtable      string   `json:"logtable"`
	User          string   `json:"user"`
	Password      string   `json:"password"`
	Timeout       int64    `json:"timeout"`
	Name          string   `json:"name"`
	SSLVerify     bool     `json:"ssl_verify"`
}

type ReDirect struct {
	RetCode int    `json:"ret_code"`
	URI     string `json:"uri"`
}
