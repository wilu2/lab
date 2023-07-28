package logs

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/base"
)

type ServiceListReq struct {
	ChannelID     int32 `form:"channel_id"`
	ApplicationID int32 `form:"application_id"`
}

type VersionListReq struct {
	ServiceID int32 `form:"service_id"`
}

type ServiceList []base.BaseInfo

type ChannelListReq struct {
	ServiceID     int32 `form:"service_id"`
	ApplicationID int32 `form:"application_id"`
	UserID        int32 `form:"user_id"`
}

type ChannelList []base.BaseInfo

type ApplicationListReq struct {
	ServiceID int32 `form:"service_id"`
	ChannelID int32 `form:"channel_id"`
}

type ApplicationList []base.BaseInfo

type DateRange struct {
	BeginDate int64  `json:"begin_date"`
	EndDate   int64  `json:"end_date"`
	Kind      string `json:"kind"`
}

type GraphReq struct {
	ServiceID     int     `json:"service_id"`
	ChannelID     int     `json:"channel_id"`
	ApplicationID []int32 `json:"application_id"`
	DateRange
}

type QueryReq struct {
	Page          int      `form:"page"`
	PageSize      int      `form:"page_size"`
	ServiceID     int      `json:"service_id"`
	ChannelID     int      `json:"channel_id"`
	ApplicationID []int32  `json:"application_id"`
	RequestID     string   `json:"request_id"`
	Version       []string `json:"version"`
	Status        []int32  `json:"status"`
	DateRange
}

type LogItem struct {
	RequestID          string  `json:"request_id"`
	ClientAddr         string  `json:"client_addr"`
	IsoTime            string  `json:"iso_time"`
	Timestamp          int64   `json:"timestamp"`
	TimeCost           float64 `json:"time_cost"`
	RequestLength      int32   `json:"request_length"`
	Connection         string  `json:"connection"`
	ConnectionRequests string  `json:"connection_requests"`
	Uri                string  `json:"uri"`
	OriRequest         string  `json:"ori_request"`
	QueryString        string  `json:"query_string"`
	Status             int32   `json:"status"`
	BytesSent          int32   `json:"bytes_sent"`
	Referer            string  `json:"referer"`
	UserAgent          string  `json:"user_agent"`
	ForwardedFor       string  `json:"forwarded_for"`
	Host               string  `json:"host"`
	Node               string  `json:"node"`
	Upstream           string  `json:"upstream"`
}

type Logs struct {
	Items []LogItem `json:"items"`
	Count int       `json:"count"`
}

type BodyLoggerReq struct {
	RequestID string `form:"request_id"`
}

type ReqBody struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	ReqBody []byte `json:"req_body"`
}

type RespBody struct {
	RespBody string `json:"resp_body"`
}
