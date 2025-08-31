package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	GetEodPriceBySymbolURL       string `env:"GET_EOD_PRICE_BY_SYMBOL_URL"`
	GetEodPriceBySecurityTypeURL string `env:"GET_EOD_PRICE_BY_SECURITY_TYPE_URL"`
	GetFinancialDataBySymbolURL  string `env:"GET_FINANCIAL_DATA_BY_SYMBOL_URL"`
	SetApiKey                    string `env:"SET_API_KEY"`
}

func C() AppConfig {
	// Load .env file if it exists (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file:", err)
	}

	// Parse environment variables into the config struct
	cfg := AppConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Unable to parse environment variables: %v", err)
	}

	// Validate required fields
	if cfg.GetEodPriceBySymbolURL == "" {
		log.Fatalf("GET_EOD_PRICE_BY_SYMBOL_URL environment variable is required")
	}
	if cfg.SetApiKey == "" {
		log.Fatalf("SET_API_KEY environment variable is required")
	}

	return cfg
}
