package entity

type GenerateSET100ReportWithTargetReq struct {
	StartDate        string   `json:"startDate"`
	EndDate          string   `json:"endDate"`
	Symbols          []string `json:"symbols"`
	TargetPercentage float64  `json:"targetPercentage"` // 21.1
	RangeOfTolerance float64  `json:"rangeOfTolerance"`
}

type GenerateSET100ReportWithTargetRespData struct {
	Symbol           string  `json:"symbol"`
	Low              float64 `json:"low"`
	LowDate          string  `json:"lowDate"`
	High             float64 `json:"high"`
	HighDate         string  `json:"highDate"`
	Target           float64 `json:"target"`
	LatestTargetDate string  `json:"latestTargetDate"`
}

type GenerateSET100ReportWithTargetResp struct {
	Code    int                                      `json:"code"`
	Message string                                   `json:"message"`
	Data    []GenerateSET100ReportWithTargetRespData `json:"data"`
}
