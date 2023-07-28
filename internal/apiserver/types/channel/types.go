package channel

type EditInfo struct {
	ID         int    `json:"id"`
	Creator    string `json:"creator"`
	CreateTime int64  `json:"create_time"`
	LastEditor string `json:"last_editor"`
	UpdateTime int64  `json:"update_time"`
}

type ChannelDef struct {
	Name string `json:"name"`
}
type Channel struct {
	EditInfo
	ChannelDef
}

type UpdateChannelReq struct {
	ID int `uri:"id"`
	ChannelDef
}

type ListReq struct {
	ID       int    `form:"id"`
	Name     string `form:"name"`      //引擎搜索关键字
	Page     int    `form:"page"`      //页数
	PageSize int    `form:"page_size"` //页大小
}

type ChannelList struct {
	Items []Channel `json:"items"`
	Count uint64    `json:"count"`
}

type IdPathParam struct {
	ID int32 `uri:"id"`
}
