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
	yesterday := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	fmt.Println("T-1 date: ", yesterday)
	yesterdayEodPriceReq := entity.BuildGetEodPriceBySecurityTypeReq(req.SecurityType, yesterday, "Y")
	reportEODT1, err := s.adapter.GetEodPriceBySecurityType(ctx, yesterdayEodPriceReq)
	if err != nil {
		return nil, err
	}
	fmt.Println("T-1 report lens: ", len(reportEODT1))
	now := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	fmt.Println("T date: ", now)
	todayEodPriceReq := entity.BuildGetEodPriceBySecurityTypeReq(req.SecurityType, now, "Y")
	reportEODT, err := s.adapter.GetEodPriceBySecurityType(ctx, todayEodPriceReq)
	if err != nil {
		return nil, err
	}
	fmt.Println("T report lens: ", len(reportEODT))
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

	return &entity.SetTopGainerResp{
		Code:    "0000",
		Message: "success",
		Data:    gainer, // Limit the number of top gainers
	}, nil
}
