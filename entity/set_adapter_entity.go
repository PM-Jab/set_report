package entity

type EodPriceBySymbol struct {
	Date              string  `json:"date"`
	Symbol            string  `json:"symbol"`
	SecurityType      string  `json:"securityType"`
	AdjustedPriceFlag string  `json:"adjustedPriceFlag"`
	Prior             float64 `json:"prior"`
	Open              float64 `json:"open"`
	High              float64 `json:"high"`
	Low               float64 `json:"low"`
	Close             float64 `json:"close"`
	Average           float64 `json:"average"`
	AomVolume         float64 `json:"aomVolume"`
	AomValue          float64 `json:"aomValue"`
	TrVolume          float64 `json:"trVolume"`
	TrValue           float64 `json:"trValue"`
	TotalVolume       float64 `json:"totalVolume"`
	TotalValue        float64 `json:"totalValue"`
	Pe                float64 `json:"pe"`
	Pbv               float64 `json:"pbv"`
	Bvps              float64 `json:"bvps"`
	DividendYield     float64 `json:"dividendYield"`
	MarketCap         float64 `json:"marketCap"`
	VolumeTurnover    float64 `json:"volumeTurnover"`
}

type GetEodPriceBySymbolReq struct {
	Symbol            string `json:"symbol"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
	AdjustedPriceFlag string `json:"adjustedPriceFlag"`
}

type GetEodPriceBySymbolResp struct {
	Code     int                `json:"code"`
	Message  string             `json:"message"`
	Response []EodPriceBySymbol `json:"response"`
}

func BuildGetEodPriceBySymbolReq(symbol, startDate, endDate, adjustedPriceFlag string) GetEodPriceBySymbolReq {
	return GetEodPriceBySymbolReq{
		Symbol:            symbol,
		StartDate:         startDate,
		EndDate:           endDate,
		AdjustedPriceFlag: adjustedPriceFlag,
	}
}
