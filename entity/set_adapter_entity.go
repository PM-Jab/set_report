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

type GetEodPriceBySecurityTypeReq struct {
	SecurityType      string `json:"securityType"`
	Date              string `json:"date"`
	AdjustedPriceFlag string `json:"adjustedPriceFlag"`
}

type GetEodPriceBySecurityTypeResp struct {
	Code     int                `json:"code"`
	Message  string             `json:"message"`
	Response []EodPriceBySymbol `json:"response"`
}

func BuildGetEodPriceBySecurityTypeReq(securityType, date, adjustedPriceFlag string) GetEodPriceBySecurityTypeReq {
	return GetEodPriceBySecurityTypeReq{
		SecurityType:      securityType,
		Date:              date,
		AdjustedPriceFlag: adjustedPriceFlag,
	}
}

type GetFinancialDataBySymbolReq struct {
	Symbol       string `json:"symbol"`
	StartYear    string `json:"startYear"`
	StartQuarter string `json:"startQuarter"`
	EndYear      string `json:"endYear"`
	EndQuarter   string `json:"endQuarter"`
}

type SetFinancialData struct {
	Symbol                 string  `json:"symbol"`
	Year                   string  `json:"year"`
	Quarter                string  `json:"quarter"`
	FinancialStatementType string  `json:"financialStatementType"`
	DateAsof               string  `json:"dateAsof"`
	AccountPeriod          string  `json:"accountPeriod"`
	TotalAssets            float64 `json:"totalAssets"`
	TotalLiabilities       float64 `json:"totalLiabilities"`
	PaidupShareCapital     float64 `json:"paidupShareCapital"`
	ShareholderEquity      float64 `json:"shareholderEquity"`
	TotalEquity            float64 `json:"totalEquity"`
	TotalRevenueQuarter    float64 `json:"totalRevenueQuarter"`
	TotalRevenueAccum      float64 `json:"totalRevenueAccum"`
	TotalExpensesQuarter   float64 `json:"totalExpensesQuarter"`
	TotalExpensesAccum     float64 `json:"totalExpensesAccum"`
	EbitQuarter            float64 `json:"ebitQuarter"`
	EbitAccum              float64 `json:"ebitAccum"`
	NetProfitQuarter       float64 `json:"netProfitQuarter"`
	NetProfitAccum         float64 `json:"netProfitAccum"`
	EpsQuarter             float64 `json:"epsQuarter"`
	EpsAccum               float64 `json:"epsAccum"`
	OperatingCashFlow      float64 `json:"operatingCashFlow"`
	InvestingCashFlow      float64 `json:"investingCashFlow"`
	FinancingCashFlow      float64 `json:"financingCashFlow"`
	Roe                    float64 `json:"roe"`
	Roa                    float64 `json:"roa"`
	NetProfitMarginQuarter float64 `json:"netProfitMarginQuarter"`
	NetProfitMarginAccum   float64 `json:"netProfitMarginAccum"`
	De                     float64 `json:"de"`
	FixedAssetTurnover     float64 `json:"fixedAssetTurnover"`
	TotalAssetTurnover     float64 `json:"totalAssetTurnover"`
}

type GetFinancialDataBySymbolResp struct {
	Code     int                `json:"code"`
	Message  string             `json:"message"`
	Response []SetFinancialData `json:"response"`
}

func BuildGetFinancialDataBySymbolReq(symbol, startYear, startQuarter, endYear, endQuarter string) GetFinancialDataBySymbolReq {
	return GetFinancialDataBySymbolReq{
		Symbol:       symbol,
		StartYear:    startYear,
		StartQuarter: startQuarter,
		EndYear:      endYear,
		EndQuarter:   endQuarter,
	}
}
