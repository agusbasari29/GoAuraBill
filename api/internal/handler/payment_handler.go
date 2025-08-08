package handler

import (
	"net/http"
	"strconv"

	"github.com/agusbasari29/GoAuraBill/internal/service"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

type CreateChargeRequest struct {
	Method string `json:"method" binding:"required"` // e.g., "BRIVA", "QRIS"
}

func (h *PaymentHandler) CreateCharge(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("transaction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID transaksi tidak valid"})
		return
	}

	var req CreateChargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chargeResponse, err := h.service.CreateTripayCharge(uint(id), req.Method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chargeResponse)
}

func (h *PaymentHandler) HandleNotification(c *gin.Context) {
	signature := c.GetHeader("X-Callback-Signature")
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.HandleTripayCallback(payload, signature); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi berhasil diproses"})
}