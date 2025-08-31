package service

import (
	"context"
	"fmt"
	"math"
	"set-report/entity"
	"sort"
	"strconv"
	"sync"
	"time"
)

// roundToTwoDecimals rounds a float64 to 2 decimal places
func roundToTwoDecimals(val float64) float64 {
	return math.Round(val*100) / 100
}

func (s *service) ScoringStockBySymbols(ctx context.Context, req entity.ScoringStockBySymbolsReq) (entity.ScoringStockBySymbolsResp, error) {
	fmt.Println("req: ", req)
	startQuarter := "4"
	endQuarter := ""

	switch month := time.Now().Month(); month {
	case 1, 2, 3:
		endQuarter = "1"
	case 4, 5, 6:
		endQuarter = "2"
	case 7, 8, 9:
		endQuarter = "3"
	case 10, 11, 12:
		endQuarter = "4"
	}

	// Create a cancellable context to stop all goroutines on error
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Ensure resources are cleaned up

	var mutex sync.Mutex
	var wg sync.WaitGroup
	var responses entity.ScoringStockBySymbolsResp
	var firstError error

	// Create a buffered channel to limit concurrency to 20 goroutines
	sem := make(chan struct{}, 20)

	for _, symbol := range req.Symbols {
		wg.Add(1)
		sem <- struct{}{} // acquire a slot
		go func(symbol string) {
			defer wg.Done()
			defer func() { <-sem }() // release the slot

			// Check if context is cancelled before proceeding
			select {
			case <-ctx.Done():
				return // Stop if context is cancelled
			default:
			}
			endYearInt, err := strconv.Atoi(req.EndYear)
			if err != nil {
				fmt.Println("Error converting EndYear to int for symbol:", symbol, "Error:", err)
				mutex.Lock()
				if firstError == nil {
					firstError = fmt.Errorf("error converting EndYear to int for symbol %s: %v", symbol, err)
				}
				mutex.Unlock()
				return
			}
			startYear := endYearInt - 5

			financialReq := entity.BuildGetFinancialDataBySymbolReq(symbol, strconv.Itoa(startYear), startQuarter, req.EndYear, endQuarter)

			data, err := s.adapter.GetFinancialDataBySymbol(ctx, financialReq)
			if err != nil {
				fmt.Println("Error fetching financial data for symbol:", symbol, "Error:", err)
				mutex.Lock()
				if firstError == nil { // Only set the first error
					firstError = fmt.Errorf("error fetching financial data for symbol %s: %v", symbol, err)
					//cancel() // Cancel context to stop all other goroutines
				}
				mutex.Unlock()
				return
			}

			// Check for cancellation again before heavy processing
			select {
			case <-ctx.Done():
				return
			default:
			}

			var financialData []entity.FinancialData
			for _, record := range data {
				// Round to 2 decimal places
				if record.Quarter == "4" {
					financialData = append(financialData, entity.FinancialData{
						Roe:                    roundToTwoDecimals(record.Roe),
						Roa:                    roundToTwoDecimals(record.Roa),
						NetProfitMarginQuarter: roundToTwoDecimals(record.NetProfitMarginQuarter),
						NetProfitMarginAccum:   roundToTwoDecimals(record.NetProfitMarginAccum),
						De:                     roundToTwoDecimals(record.De),
						OperatingCashFlow:      roundToTwoDecimals(record.OperatingCashFlow),
						TotalRevenueAccum:      roundToTwoDecimals(record.TotalRevenueAccum),
						NetProfitAccum:         roundToTwoDecimals(record.NetProfitAccum),
						EpsQuarter:             roundToTwoDecimals(record.EpsQuarter),
						EpsAccum:               roundToTwoDecimals(record.EpsAccum),
						Year:                   record.Year,
						Quarter:                record.Quarter,
					})
				}
			}

			if isFundamentallyStrongWithHistory(financialData) {
				fmt.Println("The stock", symbol, "is fundamentally strong based on the criteria.")
			} else {
				fmt.Println("The stock", symbol, "does not meet the criteria for being fundamentally strong.")
			}

			score := calculateFundamentalScore(financialData)
			fmt.Println("Fundamental Score for", symbol, "is:", score)

			// Final check before updating shared data
			select {
			case <-ctx.Done():
				return
			default:
			}

			mutex.Lock()
			responses.Data = append(responses.Data, entity.ScoringStockBySymbolsData{
				Symbol:           symbol,
				FundamentalScore: score,
			})
			mutex.Unlock()
		}(symbol)
	}

	wg.Wait()

	// Check if any error occurred
	if firstError != nil {
		fmt.Println("The first error (if any):", firstError)
		// return entity.ScoringStockBySymbolsResp{}, firstError
	}

	return responses, nil
}

func isFundamentallyStrongWithHistory(financialData []entity.FinancialData) bool {
	if len(financialData) != 5 {
		fmt.Println("Insufficient data: expected 5 years, got", len(financialData))
		return false // Requires 5 years of data
	}

	// --- Filter 1: Profitability Checks ---
	positiveProfitYears := 0
	for _, data := range financialData {
		if data.NetProfitAccum > 0 {
			positiveProfitYears++
		}
	}
	if positiveProfitYears < 4 {
		return false
	}
	if financialData[4].EpsAccum <= financialData[0].EpsAccum {
		return false
	}

	// --- Filter 2: Financial Health Checks ---
	for _, data := range financialData {
		if data.De >= 1.2 {
			return false
		}
	}

	// --- Filter 3: Efficiency Checks (calculating 5-year averages) ---
	var totalRoe, totalRoa, totalMargin float64
	for _, data := range financialData {
		totalRoe += data.Roe
		totalRoa += data.Roa
		totalMargin += data.NetProfitMarginAccum
	}
	if (totalRoe/5) <= 0.15 || (totalRoa/5) <= 0.10 || (totalMargin/5) <= 0.10 {
		return false
	}

	// --- Filter 4: Cash Flow Checks ---
	positiveCashFlowYears := 0
	for _, data := range financialData {
		if data.OperatingCashFlow > 0 {
			positiveCashFlowYears++
		}
	}
	if positiveCashFlowYears < 4 {
		return false
	}

	// --- Filter 5: Growth Checks ---
	if financialData[4].TotalRevenueAccum <= financialData[0].TotalRevenueAccum {
		return false
	}

	return true
}

func calculateFundamentalScore(financialData []entity.FinancialData) int {
	if len(financialData) != 5 {
		return 0 // Requires exactly 5 years of data
	}

	var score int = 0

	// --- 1. Profitability Scoring (Max 35 points) ---

	// Profit Consistency (15 points)
	profitConsistentYears := 0
	for _, data := range financialData {
		if data.NetProfitAccum > 0 {
			profitConsistentYears++
		}
	}
	score += profitConsistentYears * 3 // 3 points per year

	// EPS Growth (10 points)
	year1EPS := financialData[0].EpsAccum
	year5EPS := financialData[4].EpsAccum
	if year1EPS > 0 && year5EPS > (year1EPS*1.25) { // Check for > 25% growth
		score += 10
	} else if year1EPS > 0 && year5EPS > 0 {
		score += 5
	}

	// Revenue Growth (10 points)
	// Simple check: Year 5 revenue must be greater than Year 1
	if financialData[4].TotalRevenueAccum > financialData[0].TotalRevenueAccum {
		score += 10
	}

	// --- 2. Financial Health Scoring (Max 30 points) ---

	var totalDe float64
	for _, data := range financialData {
		totalDe += data.De
	}
	avgDe := totalDe / 5

	// Debt-to-Equity (DE) Ratio (20 points)
	if avgDe < 0.5 {
		score += 20
	} else if avgDe <= 1.0 {
		score += 10
	}

	// Debt Trend (10 points)
	// Check if debt has decreased or remained stable
	if financialData[4].De <= financialData[0].De {
		score += 10
	}

	// --- 3. Efficiency & Margins Scoring (Max 25 points) ---

	var totalRoe, totalMargin float64
	for _, data := range financialData {
		totalRoe += data.Roe
		totalMargin += data.NetProfitMarginAccum
	}
	avgRoe := totalRoe / 5
	avgMargin := totalMargin / 5

	// Return on Equity (ROE) (15 points)
	if avgRoe > 0.20 { // 20%
		score += 15
	} else if avgRoe >= 0.15 { // 15%
		score += 10
	} else if avgRoe >= 0.10 { // 10%
		score += 5
	}

	// Net Profit Margin (10 points)
	if avgMargin > 0.15 { // 15%
		score += 10
	} else if avgMargin >= 0.10 { // 10%
		score += 5
	}

	// --- 4. Cash Flow Quality Scoring (Max 10 points) ---
	for _, data := range financialData {
		if data.OperatingCashFlow > 0 {
			score += 2 // 2 points per year
		}
	}

	return score
}

func (s *service) FindFundamentallyStrongStockFromAllSET(ctx context.Context, req entity.FindFundamentallyStrongStockFromAllSETReq) (entity.ScoringStockBySymbolsResp, error) {
	symbols, err := s.GetAllSymbol(ctx)
	if err != nil {
		return entity.ScoringStockBySymbolsResp{}, fmt.Errorf("error fetching all symbols: %v", err)
	}

	scoringReq := entity.ScoringStockBySymbolsReq{
		Symbols:   symbols,
		StartYear: req.StartYear,
		EndYear:   req.EndYear,
	}

	ranking, err := s.ScoringStockBySymbols(ctx, scoringReq)
	if err != nil {
		return entity.ScoringStockBySymbolsResp{}, fmt.Errorf("error scoring stocks: %v", err)
	}

	// sort by FundamentalScore descending
	sort.Slice(ranking.Data, func(i, j int) bool {
		return ranking.Data[i].FundamentalScore > ranking.Data[j].FundamentalScore
	})

	return ranking, nil
}

func calculateOwnerEarnings(netIncome, depAmort, capEx float64) float64 {
	return netIncome + depAmort - capEx
}

func calculatePresentValue(futureValue, discountRate float64, year int) float64 {
	if year == 0 {
		return futureValue // Current year's value
	}
	return futureValue / math.Pow(1+discountRate, float64(year))
}

func calculateFairStockPrice(financials entity.CompanyFinancials, inputs entity.ValuationInputs) (float64, float64, error) {
	// --- Input Validation ---
	if len(financials.NetIncome) == 0 ||
		len(financials.NetIncome) != len(financials.DepreciationAmortization) ||
		len(financials.NetIncome) != len(financials.CapitalExpenditures) {
		return 0, 0, fmt.Errorf("invalid historical financial data provided")
	}
	if financials.SharesOutstanding <= 0 {
		return 0, 0, fmt.Errorf("shares outstanding must be greater than zero")
	}
	if inputs.ExplicitForecastPeriod <= 0 {
		return 0, 0, fmt.Errorf("explicit forecast period must be greater than zero")
	}
	if inputs.DiscountRate <= 0 {
		return 0, 0, fmt.Errorf("discount rate must be greater than zero")
	}
	if inputs.TerminalGrowthRate >= inputs.DiscountRate {
		return 0, 0, fmt.Errorf("terminal growth rate must be less than discount rate for stable growth model")
	}
	if inputs.MarginOfSafetyPercent < 0 || inputs.MarginOfSafetyPercent >= 1 {
		return 0, 0, fmt.Errorf("margin of safety percentage must be between 0 and 1 (exclusive)")
	}

	// 1. Calculate Historical Owner's Earnings
	historicalOwnerEarnings := make([]float64, len(financials.NetIncome))
	for i := range financials.NetIncome {
		historicalOwnerEarnings[i] = calculateOwnerEarnings(
			financials.NetIncome[i],
			financials.DepreciationAmortization[i],
			financials.CapitalExpenditures[i],
		)
	}

	// Get the most recent historical owner's earnings as the base for projection.
	lastHistoricalOE := historicalOwnerEarnings[len(historicalOwnerEarnings)-1]

	// 2. Project Future Owner's Earnings (Explicit Forecast Period)
	projectedOwnerEarnings := make([]float64, inputs.ExplicitForecastPeriod)
	presentValuesOfProjectedEarnings := make([]float64, inputs.ExplicitForecastPeriod)
	totalPVExplicitEarnings := 0.0

	for year := 1; year <= inputs.ExplicitForecastPeriod; year++ {
		// Compound growth from the last historical OE
		projectedOE := lastHistoricalOE * math.Pow(1+inputs.AnnualGrowthRate, float64(year))
		projectedOwnerEarnings[year-1] = projectedOE

		pvOE := calculatePresentValue(projectedOE, inputs.DiscountRate, year)
		presentValuesOfProjectedEarnings[year-1] = pvOE
		totalPVExplicitEarnings += pvOE
	}

	// 3. Calculate Terminal Value
	// Using the Gordon Growth Model (Perpetuity Growth Model)
	terminalYearOE := projectedOwnerEarnings[inputs.ExplicitForecastPeriod-1]

	// Ensure terminal growth rate is less than discount rate for this model
	if inputs.DiscountRate <= inputs.TerminalGrowthRate {
		return 0, 0, fmt.Errorf("discount rate must be greater than terminal growth rate for Gordon Growth Model")
	}

	terminalValue := terminalYearOE * (1 + inputs.TerminalGrowthRate) / (inputs.DiscountRate - inputs.TerminalGrowthRate)

	// Calculate Present Value of Terminal Value
	pvTerminalValue := calculatePresentValue(terminalValue, inputs.DiscountRate, inputs.ExplicitForecastPeriod)

	// 4. Calculate Total Intrinsic Value
	totalIntrinsicValue := totalPVExplicitEarnings + pvTerminalValue

	// 5. Calculate Estimated Fair Value Per Share
	fairValuePerShare := totalIntrinsicValue / financials.SharesOutstanding

	// 6. Calculate Recommended Maximum Buy Price (with Margin of Safety)
	recommendedMaxBuyPrice := fairValuePerShare * (1 - inputs.MarginOfSafetyPercent)

	return fairValuePerShare, recommendedMaxBuyPrice, nil
}
