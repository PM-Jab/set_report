package entity

type ScoringStockBySymbolsReq struct {
	Symbols   []string `json:"symbols"`
	StartYear string   `json:"startYear"`
	EndYear   string   `json:"endYear"`
}

type FinancialData struct {
	Roe                    float64
	Roa                    float64
	NetProfitMarginQuarter float64
	NetProfitMarginAccum   float64
	De                     float64
	OperatingCashFlow      float64
	TotalRevenueAccum      float64
	NetProfitAccum         float64
	EpsQuarter             float64
	EpsAccum               float64
	Year                   string
	Quarter                string
}

type ScoringStockBySymbolsData struct {
	Symbol           string `json:"symbol"`
	FundamentalScore int    `json:"fundamentalScore"`
}

type ScoringStockBySymbolsResp struct {
	Data []ScoringStockBySymbolsData `json:"data"`
}

type FindFundamentallyStrongStockFromAllSETReq struct {
	StartYear string `json:"startYear"`
	EndYear   string `json:"endYear"`
}

type CompanyFinancials struct {
	NetIncome                []float64 // Historical Net Income for multiple years
	DepreciationAmortization []float64 // Historical Depreciation & Amortization
	CapitalExpenditures      []float64 // Historical Capital Expenditures
	SharesOutstanding        float64   // Current Number of Shares Outstanding
}

// ValuationInputs holds the user-provided estimates for valuation.
type ValuationInputs struct {
	ExplicitForecastPeriod int     // Number of years for explicit cash flow projection
	AnnualGrowthRate       float64 // Annual growth rate of Owner's Earnings (e.g., 0.05 for 5%)
	DiscountRate           float64 // Discount Rate (e.g., 0.10 for 10%)
	TerminalGrowthRate     float64 // Sustainable growth rate after explicit forecast (e.g., 0.02 for 2%)
	MarginOfSafetyPercent  float64 // Margin of safety (e.g., 0.20 for 20%)
}
