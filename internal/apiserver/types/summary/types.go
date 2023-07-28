package summary

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/graph"
)

type DetailCount struct {
	Current int64 `json:"current"`
	Stopped int64 `json:"stopped"`
	Total   int64 `json:"total"`
}

type Summary struct {
	ServiceCount     DetailCount `json:"service_count"`
	ChannelCount     DetailCount `json:"channel_count"`
	ApplicationCount DetailCount `json:"application_count"`
	RequestCount     int64       `json:"request_count"`
}

type SummaryReq struct {
	BeginDate int64  `json:"begin_date"`
	EndDate   int64  `json:"end_date"`
	Type      string `json:"type"`
}

type SummaryGraph graph.PieChartInfo
