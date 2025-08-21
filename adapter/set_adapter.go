package adapter

import (
	"context"
	"errors"
	"net/http"
	"set-report/config"
	"set-report/entity"
	"set-report/httpclient.go"
)

type SetAdapter interface {
	GetEodPriceBySymbol(ctx context.Context, req entity.GetEodPriceBySymbolReq) ([]entity.EodPriceBySymbol, error)
}

func NewSetAdapter(cfg config.AppConfig, client *http.Client) SetAdapter {
	return &setAdapter{
		cfg:    cfg,
		client: client,
	}
}

type setAdapter struct {
	cfg    config.AppConfig
	client *http.Client
}

// GetEodPriceBySymbol implements SetAdapter.
func (s *setAdapter) GetEodPriceBySymbol(ctx context.Context, req entity.GetEodPriceBySymbolReq) ([]entity.EodPriceBySymbol, error) {
	url := s.cfg.GetEodPriceBySymbolURL + "?symbol=" + req.Symbol + "&startDate=" + req.StartDate + "&endDate=" + req.EndDate + "&adjustedPriceFlag=" + req.AdjustedPriceFlag
	resp, err := httpclient.Get[[]entity.EodPriceBySymbol](ctx, s.client, url, &s.cfg.SetApiKey, nil)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, err
	}

	if resp.Response == nil {
		return nil, errors.New("no data found")
	}

	return resp.Response, nil
}

func (s *setAdapter) GetEodPriceBySecurityType(ctx context.Context, req entity.GetEodPriceBySecurityTypeReq) ([]entity.EodPriceBySymbol, error) {
	url := s.cfg.GetEodPriceBySecurityTypeURL + "?securityType=" + req.SecurityType + "&date=" + req.Date + "&adjustedPriceFlag=" + req.AdjustedPriceFlag
	resp, err := httpclient.Get[[]entity.EodPriceBySymbol](ctx, s.client, url, &s.cfg.SetApiKey, nil)

	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, err
	}

	if resp.Response == nil {
		return nil, errors.New("no data found")
	}

	return resp.Response, nil
}
