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
func (s *service) GenerateTargetReportWithTargetBySymbol(ctx context.Context, req entity.GenerateTargetReportWithTargetBySymbolReq) (*entity.GenerateSETReportWithTargetResp, error) {

	var targetReportData []entity.TargetReportData

	for _, symbol := range req.Symbols {
		eodPriceReq := entity.BuildGetEodPriceBySymbolReq(symbol, req.StartDate, req.EndDate, "Y")
		eodPrices, err := s.adapter.GetEodPriceBySymbol(ctx, eodPriceReq)
		if err != nil {
			fmt.Println("Error fetching EOD prices for symbol:", symbol, "Error:", err)
			continue
		}

		if len(eodPrices) == 0 {
			continue
		}

		low := eodPrices[0].Low
		high := eodPrices[0].High
		lowDate := eodPrices[0].Date
		highDate := eodPrices[0].Date
		currentPrice := eodPrices[len(eodPrices)-1].Close

		for _, dayReport := range eodPrices {
			if dayReport.Low < low {
				low = roundToTwoDecimals(dayReport.Low)
				lowDate = dayReport.Date
			}
			if dayReport.High > high {
				high = roundToTwoDecimals(dayReport.High)
				highDate = dayReport.Date
			}
		}

		target := roundToTwoDecimals(high * req.TargetPercentage / 100)
		latestTargetDate := "" // Reset for next symbol

		// Find the latest date where Close is below target (iterate in reverse)
		for i := len(eodPrices) - 1; i >= 0; i-- {
			dayReport := eodPrices[i]
			if dayReport.Close < target {
				latestTargetDate = dayReport.Date
				break
			}
		}

		if latestTargetDate == "" {
			fmt.Println("No date found where price is below target for symbol:", symbol)
			latestTargetDate = "N/A"
		}

		targetReportData = append(targetReportData, entity.TargetReportData{
			Symbol:           symbol,
			Low:              low,
			LowDate:          lowDate,
			High:             high,
			HighDate:         highDate,
			Target:           target,
			LatestTargetDate: latestTargetDate,
			CurrentPrice:     currentPrice,
		})
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

func (s *service) GenerateTargetReportWithTargetAllSymbol(ctx context.Context, req entity.GenerateTargetReportWithTargetAllSymbolReq) (*entity.GenerateSETReportWithTargetResp, error) {

	symbols, err := s.GetAllSymbol(ctx)
	if err != nil {
		return nil, err
	}

	// batch process 100 symbols at a time with proper synchronization
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var reportTarget []entity.TargetReportData

	batchSize := 50

	for i := 0; i < len(symbols); i += batchSize {
		end := min(i+batchSize, len(symbols))
		batch := symbols[i:end]

		wg.Add(1)
		go func(symbolBatch []string) {
			defer wg.Done()

			resp, err := s.GenerateTargetReportWithTargetBySymbol(ctx, entity.GenerateTargetReportWithTargetBySymbolReq{
				Symbols:          symbolBatch,
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
			reportTarget = append(reportTarget, resp.Data...)
			mutex.Unlock()

			fmt.Println("Processed batch with", len(symbolBatch), "symbols")
		}(batch)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	filename := GenerateTargetCSVreport(reportTarget)
	fmt.Println("CSV report generated:", filename)

	return &entity.GenerateSETReportWithTargetResp{
		Code:    "0000",
		Message: "success",
		Data:    reportTarget,
	}, nil
}

func (s *service) GenerateTargetReportWithTargetAllSymbolWithLimit(ctx context.Context, req entity.GenerateTargetReportWithTargetAllSymbolWithLimitReq) (*entity.GenerateSETReportWithTargetWithLimitResp, error) {

	symbols, err := s.GetAllSymbol(ctx)
	if err != nil {
		return nil, err
	}

	if req.Limit > 100 {
		req.Limit = 100
	}

	symbols = symbols[req.Limit*(req.Page-1) : min(req.Limit*req.Page, len(symbols))]

	// batch process 100 symbols at a time with proper synchronization
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var reportTarget []entity.TargetReportData

	batchSize := 10

	for i := 0; i < len(symbols); i += batchSize {
		end := min(i+batchSize, len(symbols))
		batch := symbols[i:end]

		wg.Add(1)
		go func(symbolBatch []string) {
			defer wg.Done()

			resp, err := s.GenerateTargetReportWithTargetBySymbol(ctx, entity.GenerateTargetReportWithTargetBySymbolReq{
				Symbols:          symbolBatch,
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
			reportTarget = append(reportTarget, resp.Data...)
			mutex.Unlock()

			fmt.Println("Processed batch with", len(symbolBatch), "symbols")
		}(batch)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// filename := GenerateTargetCSVreport(reportTarget)
	// fmt.Println("CSV report generated:", filename)

	return &entity.GenerateSETReportWithTargetWithLimitResp{
		Code:    "0000",
		Message: "success",
		Data: entity.GenerateSETReportWithTargetWithLimitRespData{
			TargetReport: reportTarget,
			Page:         req.Page,
			Limit:        req.Limit,
			TotalPages:   0,
			TotalItems:   0,
		},
	}, nil
}

func GenerateTargetCSVreport(listReportTarget []entity.TargetReportData) string {
	csvContent := "Symbol,Low,LowDate,High,HighDate,Target,LatestTargetDate,CurrentPrice\n"
	for _, report := range listReportTarget {
		csvContent += fmt.Sprintf("%s,%.2f,%s,%.2f,%s,%.2f,%s,%.2f\n",
			report.Symbol,
			report.Low,
			report.LowDate,
			report.High,
			report.HighDate,
			report.Target,
			report.LatestTargetDate,
			report.CurrentPrice,
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
