package service

import (
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
)

type EditInfo struct {
	ID         int32  `json:"id"`
	Creator    string `json:"creator"`
	CreateTime int64  `json:"create_time"`
	LastEditor string `json:"last_editor"`
	UpdateTime int64  `json:"update_time"`
}

type NodeVersion struct {
	Host    string `json:"host"`
	Port    int32  `json:"port"`
	Version string `json:"version"`
}

type ApiDef struct {
	Method string `json:"method"`
	Uris   string `json:"uri"`
}

type ApiSet []ApiDef

type ServiceDef struct {
	ServiceType     string             `json:"service_type"`
	UpstreamDef     apisix.UpstreamDef `json:"upstream_def"` //上游设置
	NodeMap         []NodeVersion      `json:"node_map"`
	ApiSet          ApiSet             `json:"api_set"`
	DocumentContent string             `json:"document_content"`
}

type Service struct {
	EditInfo
	ServiceType     string          `json:"service_type"`
	Upstream        apisix.Upstream `json:"upstream_def"` //上游设置
	NodeMap         []NodeVersion   `json:"node_map"`
	ApiSet          ApiSet          `json:"api_set"`
	DocumentID      int32           `json:"document_id"`
	DocumentContent string          `json:"document_content"`
}

type ListReq struct {
	ID       int    `form:"id"`
	Name     string `form:"name"`      //引擎搜索关键字
	Page     int    `form:"page"`      //页数
	PageSize int    `form:"page_size"` //页大小
}

type ServiceList struct {
	Items []Service `json:"items"`
	Count uint64    `json:"count"`
}

type ServiceUpdate struct {
	UpstreamUpdate  apisix.UpstreamDef `json:"upstream_def"` //上游设置
	ApiSet          ApiSet             `json:"api_set"`
	NodeMap         []NodeVersion      `json:"node_map"`
	DocumentContent string             `json:"document_content"`
	ServiceType     string             `json:"service_type"`
}

type UpdateService struct {
	ID int32 `uri:"id"`
	ServiceUpdate
}

type IdPathParam struct {
	ID int32 `uri:"id"`
}

type Upstream struct {
	Host string
	Port int32
}
