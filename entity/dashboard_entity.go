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
	Code    string                                   `json:"code"`
	Message string                                   `json:"message"`
	Data    []GenerateSET100ReportWithTargetRespData `json:"data"`
}

type SetTopGainerReq struct {
	SecurityType string `json:"securityType"` // e.g., "SET"
	Limit        int    `json:"limit"`        // Number of top gainers to return
}

type SetTopData struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Change float64 `json:"change"`
}

type SetTopGainerRespData struct {
	TDate     string       `json:"tDate"`
	T1Date    string       `json:"t1Date"`
	TopGainer []SetTopData `json:"topGainer"`
}

type SetTopLoserReq struct {
	SecurityType string `json:"securityType"` // e.g., "SET"
	Limit        int    `json:"limit"`        // Number of top losers to return
}

type SetTopGainerResp struct {
	Code    string               `json:"code"`
	Message string               `json:"message"`
	Data    SetTopGainerRespData `json:"data"`
}

type SetTopLoserRespData struct {
	TDate    string       `json:"tDate"`
	T1Date   string       `json:"t1Date"`
	TopLoser []SetTopData `json:"topLoser"`
}

type SetTopLoserResp struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Data    SetTopLoserRespData `json:"data"`
}
