package main

import (
	"log/slog"
	"net/http"
	"os"
	"set-report/adapter"
	"set-report/handler"
	"set-report/service"

	"set-report/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.C()
	client := &http.Client{}
	setAdapter := adapter.NewSetAdapter(cfg, client)
	svc := service.NewService(cfg, client, setAdapter)
	h := handler.NewHandler(svc)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Define your routes here
	router.POST("/set-report", h.GenerateTargetReport)

	if err := router.Run(":" + os.Getenv("PORT")); err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
	}
}
