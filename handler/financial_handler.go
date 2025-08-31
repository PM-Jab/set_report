package handler

import (
	"set-report/entity"
	"set-report/service"

	"github.com/gin-gonic/gin"
)

type FinancialHandler struct {
	svc service.Service
}

func NewFinancialHandler(svc service.Service) *FinancialHandler {
	return &FinancialHandler{
		svc: svc,
	}
}

func (h *FinancialHandler) ScoringStockBySymbols(c *gin.Context) {
	ctx := c.Request.Context()
	var body entity.ScoringStockBySymbolsReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request parameters"})
		return
	}

	resp, err := h.svc.ScoringStockBySymbols(ctx, body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}

func (h *FinancialHandler) FindFundamentallyStrongStockFromAllSET(c *gin.Context) {
	ctx := c.Request.Context()

	var body entity.FindFundamentallyStrongStockFromAllSETReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request parameters"})
		return
	}
	resp, err := h.svc.FindFundamentallyStrongStockFromAllSET(ctx, body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}
