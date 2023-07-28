package application

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/base"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/service"
	"gitlab.intsig.net/textin-gateway/pkg/apisix"
)

type EditInfo struct {
	ID         int    `json:"id"`
	Creator    string `json:"creator"`
	CreateTime int64  `json:"create_time"`
	LastEditor string `json:"last_editor"`
	UpdateTime int64  `json:"update_time"`
}

type ApplicationDef struct {
	RouteDef  apisix.Route `json:"route_def"`
	ChannelID int          `json:"channel_id"`
	ServiceID int          `json:"service_id"`
	TTL       int          `json:"ttl"`
}

type Application struct {
	EditInfo
	ApplicationDef
}

type UpdateApplicationInfo struct {
	ID int `uri:"id"`
	ApplicationDef
}

type IdPathParam struct {
	ID int `uri:"id"`
}

type ListReq struct {
	ServiceID   int32  `form:"service_id"`
	ChannelID   int32  `form:"channel_id"`
	ServiceType string `form:"service_type"`
	ID          int    `form:"id"`
	Name        string `form:"name"`
	Status      int    `form:"status"`    //引擎搜索关键字
	Page        int    `form:"page"`      //页数
	PageSize    int    `form:"page_size"` //页大小
}

type ApplicationListItem struct {
	Application
	APISet      service.ApiSet `json:"api_set"`
	DocumentID  *int32         `json:"document_id"`
	ServiceInfo base.BaseInfo  `json:"service"`
	ChannelInfo base.BaseInfo  `json:"channel"`
}

type ApplicationList struct {
	Items []ApplicationListItem `json:"items"`
	Count uint64                `json:"count"`
}

type UpdateStatusInfo struct {
	ID     int   `uri:"id"`
	Status int32 `json:"status"`
}
