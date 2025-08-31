package handler

import (
	"fmt"
	"set-report/entity"
	"set-report/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) GenerateTargetReportWithTargetBySymbol(c *gin.Context) {
	ctx := c.Request.Context()
	var body entity.GenerateTargetReportWithTargetBySymbolReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println("Starting to generate target report with body:", body)

	resp, err := h.svc.GenerateTargetReportWithTargetBySymbol(ctx, body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate report", "details": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *Handler) SetTopGainer(c *gin.Context) {
	ctx := c.Request.Context()
	var body entity.SetTopGainerReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println("Starting to get top gainers with body:", body)

	resp, err := h.svc.SetTopGainerByDay(ctx, body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get top gainers", "details": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *Handler) SetTopLoser(c *gin.Context) {
	ctx := c.Request.Context()
	var body entity.SetTopLoserReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println("Starting to get top losers with body:", body)

	resp, err := h.svc.SetTopLoserByDay(ctx, body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get top losers", "details": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *Handler) GetAllSymbol(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := h.svc.GetAllSymbol(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get all symbols", "details": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *Handler) GenerateTargetReportWithTargetAllSymbol(c *gin.Context) {
	ctx := c.Request.Context()
	var body entity.GenerateTargetReportWithTargetAllSymbolReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println("Starting to generate target report for all symbols with body:", body)

	resp, err := h.svc.GenerateTargetReportWithTargetAllSymbol(ctx, body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate report", "details": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *Handler) GenerateTargetReportWithTargetAllSymbolWithLimit(c *gin.Context) {
	ctx := c.Request.Context()
	var body entity.GenerateTargetReportWithTargetAllSymbolWithLimitReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println("Starting to generate target report for all symbols with limit body:", body)

	resp, err := h.svc.GenerateTargetReportWithTargetAllSymbolWithLimit(ctx, body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate report", "details": err.Error()})
		return
	}
	c.JSON(200, resp)
}
