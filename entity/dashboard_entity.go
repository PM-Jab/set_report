package entity

type GenerateTargetReportWithTargetBySymbolReq struct {
	StartDate        string  `json:"startDate"`
	EndDate          string  `json:"endDate"`
	Symbol           string  `json:"symbol"`
	TargetPercentage float64 `json:"targetPercentage"` // 21.1
}

type GenerateTargetReportWithTargetAllSymbolReq struct {
	StartDate        string  `json:"startDate"`
	EndDate          string  `json:"endDate"`
	TargetPercentage float64 `json:"targetPercentage"` // 21.1
}

type GenerateTargetReportWithTargetAllSymbolWithLimitReq struct {
	StartDate        string  `json:"startDate"`
	EndDate          string  `json:"endDate"`
	TargetPercentage float64 `json:"targetPercentage"` // 21.1
	Limit            int     `json:"limit"`            // Limit for the number of symbols
	Page             int     `json:"page"`             // Page number for pagination
}

type TargetReportData struct {
	Symbol           string  `json:"symbol"`
	Low              float64 `json:"low"`
	LowDate          string  `json:"lowDate"`
	High             float64 `json:"high"`
	HighDate         string  `json:"highDate"`
	Target           float64 `json:"target"`
	LatestTargetDate string  `json:"latestTargetDate"`
	CurrentPrice     float64 `json:"currentPrice"`
	CurrentDate      string  `json:"currentDate"`
}

type GenerateSETReportWithTargetResp struct {
	Code    string           `json:"code"`
	Message string           `json:"message"`
	Data    TargetReportData `json:"data"`
}

type GenerateSETReportWithAllSymbolWithTargetResp struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Data    []TargetReportData `json:"data"`
}

type GenerateSETReportWithTargetWithLimitRespData struct {
	TargetReport []TargetReportData `json:"targetReport"`
	Page         int                `json:"page"`
	Limit        int                `json:"limit"`
	TotalPages   int                `json:"totalPages"`
	TotalItems   int                `json:"totalItems"`
}

type GenerateSETReportWithTargetWithLimitResp struct {
	Code    string                                       `json:"code"`
	Message string                                       `json:"message"`
	Data    GenerateSETReportWithTargetWithLimitRespData `json:"data"`
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

type MonthlyAverageStockPriceBySymbolReq struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Symbol    string `json:"symbol"`
}

type MonthlyAverageStockPriceBySymbolResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Symbol            string `json:"symbol"`
		MonthlyStockPrice []struct {
			Date  string  `json:"date"`
			Price float64 `json:"price"`
		} `json:"monthlyStockPrice"`
	} `json:"data"`
}
