package service

import (
	"context"
	"net/http"
	"set-report/adapter"
	"set-report/config"
	"set-report/entity"
)

type Service interface {
	GenerateSETReportWithTarget(ctx context.Context, req entity.GenerateSET100ReportWithTargetReq) (*entity.GenerateSET100ReportWithTargetResp, error)
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
		Code:    0,
		Message: "success",
		Data:    targetReportData,
	}, nil
}

func (s *service) GetAllSymbols(ctx context.Context) ([]string, error) {
	return s.adapter.GetAllSymbols(ctx)
}

func IsPriceInToleranceRange(price, targetPrice, tolerance float64) bool {
	if price < targetPrice-tolerance || price > targetPrice+tolerance {
		return false
	}
	return true
}
