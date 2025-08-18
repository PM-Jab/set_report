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

func (h *Handler) GenerateTargetReport(c *gin.Context) {
	ctx := c.Request.Context()
	var body entity.GenerateSET100ReportWithTargetReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println("Starting to generate target report with body:", body)

	resp, err := h.svc.GenerateSET100ReportWithTarget(ctx, body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate report", "details": err.Error()})
		return
	}
	c.JSON(200, resp)
}
