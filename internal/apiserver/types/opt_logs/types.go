package opt_logs

// type DateRange struct {
// 	BeginDate int `json:"begin_date"`
// 	EndDate   int `json:"end_date"`
// }

type ListReq struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	UserID   int32  `form:"user_id"`
	Resource string `form:"resource"`
}

type DeleteReq struct {
	BeginDate int `json:"begin_date"`
	EndDate   int `json:"end_date"`
}

type DetailReq struct {
	ID int32 `uri:"id"`
}

type OptLogItem struct {
	ID           int32  `json:"id"`
	OptTime      int64  `json:"opt_time"`
	Operation    string `json:"operation"`
	Resource     string `json:"resource"`
	ResourceType string `json:"resource_type"`
	UserName     string `json:"user_name"`
}

type OptLogs struct {
	Items []OptLogItem `json:"items"`
	Count int          `json:"count"`
}

type OptLogsDetail struct {
	Detail string `json:"detail"`
}
