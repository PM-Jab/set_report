package service

import (
	"context"
	"fmt"
	"net/http"
	"set-report/adapter"
	"set-report/config"
	"set-report/entity"
	"sort"
	"time"
)

type Service interface {
	GenerateSETReportWithTarget(ctx context.Context, req entity.GenerateSET100ReportWithTargetReq) (*entity.GenerateSET100ReportWithTargetResp, error)
	SetTopGainerByDay(ctx context.Context, req entity.SetTopGainerReq) (*entity.SetTopGainerResp, error)
	SetTopLoserByDay(ctx context.Context, req entity.SetTopLoserReq) (*entity.SetTopLoserResp, error)
}

func NewService(cfg config.AppConfig, client *http.Client, adapter adapter.SetAdapter) Service {
	return &service{
		cfg:     cfg,
		client:  client,
		adapter: adapter,
	}
}

type service struct {
	cfg     config.AppConfig
	client  *http.Client
	adapter adapter.SetAdapter
}

// GetEodPriceBySymbol implements Service.
func (s *service) GenerateSETReportWithTarget(ctx context.Context, req entity.GenerateSET100ReportWithTargetReq) (*entity.GenerateSET100ReportWithTargetResp, error) {

	var targetReportData []entity.GenerateSET100ReportWithTargetRespData
	var low, high, target float64
	var lowDate, highDate, latestTargetDate string

	for _, symbol := range req.Symbols {
		eodPriceReq := entity.BuildGetEodPriceBySymbolReq(symbol, req.StartDate, req.EndDate, "Y")
		eodPrices, err := s.adapter.GetEodPriceBySymbol(ctx, eodPriceReq)
		if err != nil {
			return nil, err
		}

		if len(eodPrices) == 0 {
			continue
		}

		low = eodPrices[0].Low
		high = eodPrices[0].High
		lowDate = eodPrices[0].Date
		highDate = eodPrices[0].Date

		for _, dayReport := range eodPrices {
			if low == 0 || dayReport.Low < low {
				low = dayReport.Close
				lowDate = dayReport.Date
			}
			if high == 0 || dayReport.High > high {
				high = dayReport.Close
				highDate = dayReport.Date
			}
		}

		target = high * req.TargetPercentage / 100

		for _, dayReport := range eodPrices {
			if IsPriceInToleranceRange(dayReport.Close, target, req.RangeOfTolerance) {
				latestTargetDate = dayReport.Date
			}
		}

		targetReportData = append(targetReportData, entity.GenerateSET100ReportWithTargetRespData{
			Symbol:           symbol,
			Low:              low,
			LowDate:          lowDate,
			High:             high,
			HighDate:         highDate,
			Target:           target,
			LatestTargetDate: latestTargetDate,
		})

		latestTargetDate = "" // Reset for next symbol

	}

	return &entity.GenerateSET100ReportWithTargetResp{
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
