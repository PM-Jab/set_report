package service

import (
	"context"
	"fmt"
	"math"
	"set-report/entity"
	"sort"
	"sync"
	"time"
)

// GetEodPriceBySymbol implements Service.
func (s *service) GenerateTargetReportWithTargetBySymbol(ctx context.Context, req entity.GenerateTargetReportWithTargetBySymbolReq) (resp *entity.GenerateSETReportWithTargetResp, err error) {

	var targetReportData entity.TargetReportData
	symbol := req.Symbol

	eodPriceReq := entity.BuildGetEodPriceBySymbolReq(symbol, req.StartDate, req.EndDate, "Y")
	eodPrices, err := s.adapter.GetEodPriceBySymbol(ctx, eodPriceReq)
	if err != nil {
		fmt.Println("Error fetching EOD prices for symbol:", symbol, "Error:", err)
		return nil, fmt.Errorf("error fetching EOD prices for symbol: %s", symbol)
	}

	if len(eodPrices) == 0 {
		return nil, fmt.Errorf("no EOD prices found for symbol: %s", symbol)
	}

	// Initialize low, high, lowDate, highDate, currentPrice
	low := eodPrices[0].Close
	lowDate := eodPrices[0].Date
	high := 0.0
	highDate := ""
	currentPrice := eodPrices[len(eodPrices)-1].Close
	currentDate := eodPrices[len(eodPrices)-1].Date

	// Find lowest
	indexOfLowDate := 0
	for i, dayReport := range eodPrices {
		if dayReport.Close < low || (low == 0 && dayReport.Close > 0) {
			low = roundToTwoDecimals(dayReport.Close)
			lowDate = dayReport.Date
			indexOfLowDate = i
		}
	}
	// Find highest after lowDate
	for _, dayReport := range eodPrices[indexOfLowDate:] {
		if dayReport.Close > high {
			high = roundToTwoDecimals(dayReport.Close)
			highDate = dayReport.Date
		}
	}

	if high == 0 {
		fmt.Println("No high price found after low date for symbol:", symbol)
		return nil, fmt.Errorf("no high price found after low date for symbol: %s", symbol)
	}

	target := roundToTwoDecimals(((high - low) * req.TargetPercentage / 100) + low)
	latestTargetDate := "" // Reset for next symbol

	// Find the latest date after highDate where Close is below target (iterate in reverse)
	for i := len(eodPrices) - 1; i >= 0; i-- {
		dayReport := eodPrices[i]
		if dayReport.Close < target {
			latestTargetDate = dayReport.Date
			break
		}

		// If we reach the highDate, stop searching
		if dayReport.Date == highDate {
			break
		}
	}

	if latestTargetDate == "" {
		fmt.Println("No date found where price is below target for symbol:", symbol)
		latestTargetDate = "N/A"
	}

	targetReportData = entity.TargetReportData{
		Symbol:           symbol,
		Low:              low,
		LowDate:          lowDate,
		High:             high,
		HighDate:         highDate,
		Target:           target,
		LatestTargetDate: latestTargetDate,
		CurrentPrice:     currentPrice,
		CurrentDate:      currentDate,
	}

	return &entity.GenerateSETReportWithTargetResp{
		Code:    "0000",
		Message: "success",
		Data:    targetReportData,
	}, nil
}

func IsPriceInToleranceRange(price, targetPrice, tolerance float64) bool {
	if price < targetPrice-tolerance || price > targetPrice+tolerance {
		return false
	}
	return true
}

func (s *service) SetTopGainerByDay(ctx context.Context, req entity.SetTopGainerReq) (*entity.SetTopGainerResp, error) {
	t_1 := 0
	var reportEODT, reportEODT1 []entity.EodPriceBySymbol
	var tDate, t1Date string

	for {
		t := time.Now().AddDate(0, 0, -t_1).Format("2006-01-02")
		todayEodPriceReq := entity.BuildGetEodPriceBySecurityTypeReq(req.SecurityType, t, "Y")
		reportEOD, err := s.adapter.GetEodPriceBySecurityType(ctx, todayEodPriceReq)
		if err != nil {
			return nil, err
		}

		if len(reportEOD) > 0 {
			if reportEODT == nil {
				reportEODT = reportEOD
				tDate = t
			} else {
				reportEODT1 = reportEOD
				t1Date = t
				break
			}
		}

		t_1++

		if t_1 > 10 {
			fmt.Println("No more data found")
			return &entity.SetTopGainerResp{}, nil
		}
	}

	var gainer []entity.SetTopData
	for index, t1 := range reportEODT1 {
		if t1.Close == 0 {
			continue
		}
		if t1.Close < reportEODT[index].Close {
			gainer = append(gainer, entity.SetTopData{
				Symbol: t1.Symbol,
				Price:  reportEODT[index].Close,
				Change: reportEODT[index].Close/t1.Close - 1,
			})
		}
	}

	// sort by change
	sort.Slice(gainer, func(i, j int) bool {
		return gainer[i].Change > gainer[j].Change
	})

	var filteredGainers []entity.SetTopData
	for _, g := range gainer {
		if g.Price >= 1 {
			filteredGainers = append(filteredGainers, g)
		}

		if filteredGainers != nil && len(filteredGainers) >= req.Limit || g.Change*100 < 5 {
			break
		}
	}

	return &entity.SetTopGainerResp{
		Code:    "0000",
		Message: "success",
		Data: entity.SetTopGainerRespData{
			TDate:     tDate,
			T1Date:    t1Date,
			TopGainer: filteredGainers,
		},
	}, nil
}

func (s *service) SetTopLoserByDay(ctx context.Context, req entity.SetTopLoserReq) (*entity.SetTopLoserResp, error) {
	t_1 := 0
	var reportEODT, reportEODT1 []entity.EodPriceBySymbol
	var tDate, t1Date string

	for {
		t := time.Now().AddDate(0, 0, -t_1).Format("2006-01-02")
		todayEodPriceReq := entity.BuildGetEodPriceBySecurityTypeReq(req.SecurityType, t, "Y")
		reportEOD, err := s.adapter.GetEodPriceBySecurityType(ctx, todayEodPriceReq)
		if err != nil {
			return nil, err
		}

		if len(reportEOD) > 0 {
			if reportEODT == nil {
				reportEODT = reportEOD
				tDate = t
			} else {
				reportEODT1 = reportEOD
				t1Date = t
				break
			}
		}

		t_1++

		if t_1 > 10 {
			fmt.Println("No more data found")
			return &entity.SetTopLoserResp{}, nil
		}
	}

	var loser []entity.SetTopData
	for index, t1 := range reportEODT1 {
		if t1.Close == 0 {
			continue
		}
		if t1.Close > reportEODT[index].Close {
			loser = append(loser, entity.SetTopData{
				Symbol: t1.Symbol,
				Price:  reportEODT[index].Close,
				Change: 1 - reportEODT[index].Close/t1.Close,
			})
		}
	}

	// sort by change
	sort.Slice(loser, func(i, j int) bool {
		return loser[i].Change > loser[j].Change
	})

	var filteredLosers []entity.SetTopData
	for _, l := range loser {
		if l.Price >= 1 {
			filteredLosers = append(filteredLosers, l)
		}

		if filteredLosers != nil && len(filteredLosers) >= req.Limit || l.Change*100 < 5 {
			break
		}
	}

	return &entity.SetTopLoserResp{
		Code:    "0000",
		Message: "success",
		Data: entity.SetTopLoserRespData{
			TDate:    tDate,
			T1Date:   t1Date,
			TopLoser: filteredLosers,
		},
	}, nil
}

func (s *service) GetAllSymbol(ctx context.Context) ([]string, error) {
	securityType := "CS"

	for t := range 10 {
		now := time.Now().AddDate(0, 0, -t).Format("2006-01-02")
		eodPriceReq := entity.BuildGetEodPriceBySecurityTypeReq(securityType, now, "Y")
		eodPrices, err := s.adapter.GetEodPriceBySecurityType(ctx, eodPriceReq)

		if err != nil {
			return nil, err
		}
		if len(eodPrices) > 0 {
			fmt.Println("Found data for date:", now)
			// Extract symbols from eodPrices
			var symbols []string
			for _, eod := range eodPrices {
				symbols = append(symbols, eod.Symbol)
			}
			return symbols, nil
		}
	}
	fmt.Println("No data found")

	return nil, nil
}

func (s *service) GenerateTargetReportWithTargetAllSymbol(ctx context.Context, req entity.GenerateTargetReportWithTargetAllSymbolReq) (*entity.GenerateSETReportWithAllSymbolWithTargetResp, error) {

	symbols, err := s.GetAllSymbol(ctx)
	if err != nil {
		return nil, err
	}

	// batch process 100 symbols at a time with proper synchronization
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var reportTarget []entity.TargetReportData
	worker := make(chan struct{}, 30) // limit to 30 concurrent goroutines

	for _, symbol := range symbols {
		worker <- struct{}{} // acquire a worker
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			defer func() { <-worker }() // release the worker

			resp, err := s.GenerateTargetReportWithTargetBySymbol(ctx, entity.GenerateTargetReportWithTargetBySymbolReq{
				Symbol:           symbol,
				StartDate:        req.StartDate,
				EndDate:          req.EndDate,
				TargetPercentage: req.TargetPercentage,
			})
			if err != nil {
				fmt.Println("Error processing batch:", err)
				return
			}

			// Use mutex to safely append to shared slice
			mutex.Lock()
			reportTarget = append(reportTarget, resp.Data)
			mutex.Unlock()

			fmt.Println("Processed batch with", symbol, "symbols")
		}(symbol)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	filename := GenerateTargetCSVreport(reportTarget)
	fmt.Println("CSV report generated:", filename)

	return &entity.GenerateSETReportWithAllSymbolWithTargetResp{
		Code:    "0000",
		Message: "success",
		Data:    reportTarget,
	}, nil
}

func (s *service) MonthlyAverageStockPriceBySymbol(ctx context.Context, req entity.MonthlyAverageStockPriceBySymbolReq) (*entity.MonthlyAverageStockPriceBySymbolResp, error) {
	eodPriceReq := entity.BuildGetEodPriceBySymbolReq(req.Symbol, req.StartDate, req.EndDate, "Y")
	eodPrices, err := s.adapter.GetEodPriceBySymbol(ctx, eodPriceReq)
	if err != nil {
		return nil, err
	}

	if len(eodPrices) == 0 {
		return nil, fmt.Errorf("no EOD prices found for symbol: %s", req.Symbol)
	}

	// Calculate monthly average stock price
	// reset every 21 days
	monthlyTotal := make(map[int]float64)
	var month, monthlyCount int
	startMonthlyDate := []string{eodPrices[len(eodPrices)-1].Date}
	for i := len(eodPrices) - 1; i >= 0; i-- {
		eod := eodPrices[i]
		if monthlyCount < 21 {
			monthlyTotal[month] += eod.Close
		} else {
			month++
			monthlyCount = 0
			monthlyTotal[month] += eod.Close
			startMonthlyDate = append(startMonthlyDate, eod.Date)
		}
		monthlyCount++
	}

	var monthlyStockPrice []struct {
		Date   string  `json:"date"`
		Price  float64 `json:"price"`
		Change float64 `json:"change"`
	}

	var totalPositive int

	for i := 0; i <= month; i++ {
		avgPrice := roundToTwoDecimals(monthlyTotal[i] / 21)
		if math.IsNaN(avgPrice) {
			return nil, fmt.Errorf("average price is NaN for symbol: %s", req.Symbol)
		}
		var change float64
		if i > 0 {
			change = roundToTwoDecimals((monthlyStockPrice[i-1].Price - avgPrice) * 100 / avgPrice)
			if math.IsNaN(change) {
				return nil, fmt.Errorf("change is NaN for symbol: %s", req.Symbol)
			}

			monthlyStockPrice[i-1].Change = change
			if change > 0 {
				totalPositive++
			}
		}
		monthlyStockPrice = append(monthlyStockPrice, struct {
			Date   string  `json:"date"`
			Price  float64 `json:"price"`
			Change float64 `json:"change"`
		}{
			Date:   startMonthlyDate[i],
			Price:  avgPrice,
			Change: 0,
		})
	}

	return &entity.MonthlyAverageStockPriceBySymbolResp{
		Code:    "0000",
		Message: "success",
		Data: entity.MonthlyAverageStockPriceBySymbolData{
			Symbol:            req.Symbol,
			TotalPositive:     totalPositive,
			MonthlyStockPrice: monthlyStockPrice,
		},
	}, nil
}

func (s *service) MonthlyAverageStockPriceAllSymbol(ctx context.Context, req entity.MonthlyAverageStockPriceAllSymbolReq) (*entity.MonthlyAverageStockPriceAllSymbolResp, error) {
	symbols, err := s.GetAllSymbol(ctx)
	if err != nil {
		return nil, err
	}

	var mutex sync.Mutex
	var wg sync.WaitGroup
	var monthlyAvgAllSymbol []entity.MonthlyAverageStockPriceBySymbolData
	worker := make(chan struct{}, 20) // limit to 20 concurrent goroutines

	for _, symbol := range symbols {
		worker <- struct{}{} // acquire a worker
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			defer func() { <-worker }() // release the worker

			resp, err := s.MonthlyAverageStockPriceBySymbol(ctx, entity.MonthlyAverageStockPriceBySymbolReq{
				Symbol:    symbol,
				StartDate: req.StartDate,
				EndDate:   req.EndDate,
			})
			if err != nil {
				fmt.Println("Error processing batch:", err)
				return
			}

			if len(resp.Data.MonthlyStockPrice) == 0 {
				return
			}

			// Use mutex to safely append to shared slice
			mutex.Lock()
			monthlyAvgAllSymbol = append(monthlyAvgAllSymbol, resp.Data)
			mutex.Unlock()

			fmt.Println("Processed batch with", symbol, "symbols")
		}(symbol)
	}
	// Wait for all goroutines to complete
	wg.Wait()

	// generate CSV report
	filename := GenerateMonthlyAverageCSVreport(monthlyAvgAllSymbol)
	fmt.Println("CSV report generated:", filename)

	return &entity.MonthlyAverageStockPriceAllSymbolResp{
		Code:    "0000",
		Message: "success",
		Data:    monthlyAvgAllSymbol,
	}, nil
}
