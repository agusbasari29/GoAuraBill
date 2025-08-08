package handler

import (
	"net/http"

	"github.com/agusbasari29/GoAuraBill/internal/service"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service service.ReportService
}

func NewReportHandler(service service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GetRevenueReport: Endpoint untuk mendapatkan laporan pendapatan.
// Menerima query parameter 'period' (daily/monthly).
func (h *ReportHandler) GetRevenueReport(c *gin.Context) {
	period := c.DefaultQuery("period", "daily") // Default ke 'daily' jika tidak disediakan
	report, err := h.service.GetRevenueReport(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

// GetSummaryReport: Endpoint untuk mendapatkan ringkasan metrik bisnis.
func (h *ReportHandler) GetSummaryReport(c *gin.Context) {
	summary, err := h.service.GetSummaryReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}
