package graph

type LineItem struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

type LineChartInfo struct {
	XAxis  []int64    `json:"xAxis"`
	Series []LineItem `json:"series"`
}

type PieItem struct {
	Name string `json:"name"`
	Data []int  `json:"data"`
}

type PieChartInfo struct {
	XAxis  []int64   `json:"xAxis"`
	YAxis  string    `json:"yAxis"`
	Series []PieItem `json:"series"`
}

type ColumnItem struct {
	Name string `json:"name"`
	Data []int  `json:"data"`
}

type ColumnChartInfo struct {
	Category []string     `json:"category"`
	Series   []ColumnItem `json:"series"`
}
