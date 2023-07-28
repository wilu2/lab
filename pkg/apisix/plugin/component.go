package plugin

type BreakResponseHeaders struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type HealthyStatus struct {
	HTTPStatuses []int `json:"http_statuses"`
	Successes    int   `json:"successes"`
}
type UnhealthyStatus struct {
	Failures     int   `json:"failures"`
	HTTPStatuses []int `json:"http_statuses"`
}
