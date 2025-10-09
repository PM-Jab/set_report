package service

import (
	"context"
	"net/http"
	"set-report/adapter"
	"set-report/config"
	"set-report/entity"
)

type Service interface {
	GenerateTargetReportWithTargetBySymbol(ctx context.Context, req entity.GenerateTargetReportWithTargetBySymbolReq) (*entity.GenerateSETReportWithTargetResp, error)
	SetTopGainerByDay(ctx context.Context, req entity.SetTopGainerReq) (*entity.SetTopGainerResp, error)
	SetTopLoserByDay(ctx context.Context, req entity.SetTopLoserReq) (*entity.SetTopLoserResp, error)
	GetAllSymbol(ctx context.Context) ([]string, error)
	GenerateTargetReportWithTargetAllSymbol(ctx context.Context, req entity.GenerateTargetReportWithTargetAllSymbolReq) (*entity.GenerateSETReportWithAllSymbolWithTargetResp, error)
	// GenerateTargetReportWithTargetAllSymbolWithLimit(ctx context.Context, req entity.GenerateTargetReportWithTargetAllSymbolWithLimitReq) (*entity.GenerateSETReportWithTargetWithLimitResp, error)
	ScoringStockBySymbols(ctx context.Context, req entity.ScoringStockBySymbolsReq) (entity.ScoringStockBySymbolsResp, error)
	FindFundamentallyStrongStockFromAllSET(ctx context.Context, req entity.FindFundamentallyStrongStockFromAllSETReq) (entity.ScoringStockBySymbolsResp, error)
	MonthlyAverageStockPriceBySymbol(ctx context.Context, req entity.MonthlyAverageStockPriceBySymbolReq) (*entity.MonthlyAverageStockPriceBySymbolResp, error)
	MonthlyAverageStockPriceAllSymbol(ctx context.Context, req entity.MonthlyAverageStockPriceAllSymbolReq) (*entity.MonthlyAverageStockPriceAllSymbolResp, error)
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
