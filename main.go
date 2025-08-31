package main

import (
	"log/slog"
	"net/http"
	"os"
	"set-report/adapter"
	"set-report/handler"
	"set-report/service"

	"set-report/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.C()
	client := &http.Client{}
	setAdapter := adapter.NewSetAdapter(cfg, client)
	svc := service.NewService(cfg, client, setAdapter)
	fsvc := service.NewService(cfg, client, setAdapter)
	h := handler.NewHandler(svc)
	fh := handler.NewFinancialHandler(fsvc)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(cors.Default())

	// Define your routes here
	router.POST("/set-report", h.GenerateTargetReportWithTargetBySymbol)
	router.POST("/set-report/limit", h.GenerateTargetReportWithTargetAllSymbolWithLimit)
	router.POST("/set-report/all", h.GenerateTargetReportWithTargetAllSymbol)
	router.POST("/set-top-gainer", h.SetTopGainer)
	router.POST("/set-top-loser", h.SetTopLoser)
	router.GET("/set-symbols", h.GetAllSymbol)

	router.POST("/financial/scoring-all", fh.FindFundamentallyStrongStockFromAllSET)
	router.POST("/financial/scoring-by-symbols", fh.ScoringStockBySymbols)

	if err := router.Run(":" + os.Getenv("PORT")); err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
	}
}
