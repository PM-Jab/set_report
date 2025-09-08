package service

import (
	"context"
	"fmt"
	"os"
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

func GenerateTargetCSVreport(listReportTarget []entity.TargetReportData) string {
	csvContent := "Symbol,Low,LowDate,High,HighDate,Target,LatestTargetDate,CurrentPrice,CurrentDate\n"
	for _, report := range listReportTarget {
		csvContent += fmt.Sprintf("%s,%.2f,%s,%.2f,%s,%.2f,%s,%.2f,%s\n",
			report.Symbol,
			report.Low,
			report.LowDate,
			report.High,
			report.HighDate,
			report.Target,
			report.LatestTargetDate,
			report.CurrentPrice,
			report.CurrentDate,
		)
	}

	// create file name with current timestamp
	filename := fmt.Sprintf("target_report_%s.csv", time.Now().Format("20060102_150405"))
	// write to file
	err := os.WriteFile(filename, []byte(csvContent), 0644)
	if err != nil {
		fmt.Println("Error writing CSV file:", err)
		return ""
	}
	fmt.Println("CSV report generated:", filename)
	return filename
}

func (s *service) MonthlyAverageStockPriceBySymbol(ctx context.Context, req entity.MonthlyAverageStockPriceBySymbolReq) (*entity.MonthlyAverageStockPriceBySymbolResp, error) {
	eodPriceReq := entity.BuildGetEodPriceBySymbolReq(req.Symbol, req.StartDate, req.EndDate, "Y")
	eodPrices, err := s.adapter.GetEodPriceBySymbol(ctx, eodPriceReq)
	if err != nil {
		return nil, err
	}

	// Calculate monthly average stock price
	monthlyAvg := make(map[string]float64)
	for _, eod := range eodPrices {
		month := eod.Date[:7] // Get the year-month part
		monthlyAvg[month] += eod.Close
	}

	// Prepare response
	var monthlyStockPrices []struct {
		Date  string  `json:"date"`
		Price float64 `json:"price"`
	}
	for month, total := range monthlyAvg {
		monthlyStockPrices = append(monthlyStockPrices, struct {
			Date  string  `json:"date"`
			Price float64 `json:"price"`
		}{
			Date:  month,
			Price: total / float64(len(eodPrices)), // Average price
		})
	}

	return &entity.MonthlyAverageStockPriceBySymbolResp{
		Code:    "0000",
		Message: "success",
		Data: struct {
			Symbol            string `json:"symbol"`
			MonthlyStockPrice []struct {
				Date  string  `json:"date"`
				Price float64 `json:"price"`
			} `json:"monthlyStockPrice"`
		}{
			Symbol:            req.Symbol,
			MonthlyStockPrice: monthlyStockPrices,
		},
	}, nil
}
